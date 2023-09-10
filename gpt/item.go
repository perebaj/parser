package gpt

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	"github.com/sashabaranov/go-openai"
)

const itemSystemPrompt = `
	This prompt will receive a list of items from a purchase order in a unstructured text format. 
	
	An input example is:
	TEXT:
		24676   AGUA
			-   PURIFICADA CRYSTAL 2L            7894900525021 PT 6X2LT
																PT         20,00        0,00   13,690002   0,00 13,690002    0,00   0,00   0,00   0,00 13,690002    0,00   0,00   0,00    0,00   0,00   0,00   0,00
		24677   AGUA
			-   PURIFICADA CRYSTAL 510ML         7894900525007 PT 12X510ML
																PT      20,00           0,00   10,070004   0,00 10,070004    0,00   0,00   0,00   0,00 10,070004    0,00   0,00   0,00    0,00   0,00   0,00   0,00
		12285   AGUA
			-   TONICA SCHWEPPES LT 350ML        7894900300017 PT 12UN
																PT          5,00        0,00   22,509996   0,00 22,509996    0,00   0,00   0,00   0,00 22,509996    0,00   0,00   0,00    0,00   0,00   0,00   0,00
	END TEXT
	
	You will need to parse the text and extract the following fields:
		- Código
		- Descrição
		- Cód.Fáb.
		- Emb
		- Un.
		- Qt.Pedida
		- Qt. Entreg.
		- P.Compra
		- %Desc.
		- P.Líquido
		- %Frete
		- %IPI
		- %ST
		- Preço Merc.
		- %ICMS

	The schema and the output must be like this:
	[
		{
			"code": "24676",
			"description": "AGUA - PURIFICADA CRYSTAL 2L",
			"factory_code": "7894900525021",
			"emb": "PT",
			"un": "6X2LT",
			"quantity_ordered": "20,00",
			"quantity_delivered": "0,00",
			"purchase_price": "13,690002",
			"discount": "0,00",
			"liquid_price": "13,690002",
			"freight": "0,00",
			"ipi": "0,00",
			"st": "0,00",
			"price_merc": "0,00",
			"icms": "0,00",
		},
		{
			"code": "24677",
			"description": "AGUA - PURIFICADA CRYSTAL 510ML",
			"factory_code": "7894900525007",
			"emb": "PT",
			"un": "12X510ML",
			"quantity_ordered": "20,00",
			"quantity_delivered": "0,00",
			"purchase_price": "10,070004",
			"discount": "0,00",
			"liquid_price": "10,070004",
			"freight": "0,00",
			"ipi": "0,00",
			"st": "0,00",
			"price_merc": "0,00",
			"icms": "0,00",
		},
		{
			"code": "12285",
			"description": "AGUA - TONICA SCHWEPPES LT 350ML",
			"factory_code": "7894900300017",
			"emb": "PT",
			"un": "12UN",
			"quantity_ordered": "5,00",
			"quantity_delivered": "0,00",
			"purchase_price": "22,509996",
			"discount": "0,00",
			"liquid_price": "22,509996",
			"freight": "0,00",
			"ipi": "0,00",
			"st": "0,00",
			"price_merc": "0,00",
			"icms": "0,00",
		}
	]

	Do not include any fields that are not specified in the schema.
	Do not Generate any fields that are not specified in the schema.
	Do not generate any data if you don't have enough confidence in it.
`

// Item represents a purchase order item
type Item struct {
	Code              string `json:"code"`
	Description       string `json:"description"`
	FactoryCode       string `json:"factory_code"`
	Emb               string `json:"emb"`
	Un                string `json:"un"`
	QuantityOrdered   string `json:"quantity_ordered"`
	QuantityDelivered string `json:"quantity_delivered"`
	PurchasePrice     string `json:"purchase_price"`
	Discount          string `json:"discount"`
	LiquidPrice       string `json:"liquid_price"`
	Freight           string `json:"freight"`
	Ipi               string `json:"ipi"`
	St                string `json:"st"`
	PriceMerc         string `json:"price_merc"`
	Icms              string `json:"icms"`
}

// ItemParser parses the items of a purchase order
func (c *Client) ItemParser(text string) ([]Item, error) {
	usrPrompt, err := createUserPrompt(text)
	if err != nil {
		return nil, fmt.Errorf("create user prompt error %v", err)
	}

	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			Temperature: math.SmallestNonzeroFloat32,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: itemSystemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: usrPrompt,
				},
			},
		})
	if err != nil {
		return nil, fmt.Errorf("gpt response error %v", err)
	}
	respContent := resp.Choices[0].Message.Content
	var items []Item
	err = json.Unmarshal([]byte(respContent), &items)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal error %v", err)

	}
	return items, nil
}
