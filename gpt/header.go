// Package gpt package is a wrapper for the OpenAI GPT-3 API
package gpt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"math"

	"github.com/sashabaranov/go-openai"
)

const systemPrompt = `
	You will work like a unstructured text parser, I will give you a text and you will parse it according to the specified schema
	and fields. You will return a JSON object.
	extractor
	An input example:

	TEXT:
		210 - Emitir Pedido de Compra - Com saldo de pedido                                                                                      Data : 03/05/2023                                                  Página 1 de 4
		*** PEDIDO EMITIDO NA UNID. MASTER ***                      Número do Pedido:                   48740        N.VerbaBonific 0            Filial: 1
																COMPRA MERC.                                                       Nro da Negociação: 0         Tipo de carga: Paletizado
		
		Empresa:           FAKE                                                                        CNPJ/CPF:           00000000000000
		Endereço:          FAKE                      915   PARQUE DEZ                                                  Inscrição Estadual: 042133815
		Cidade:            MANAUS                          UF: AM     CEP: FAKE                                                 Telefone:           0000000000
		
		Fornecedor: FAKE S/A                                                  Codigo:  114                      CNPJ/CPF:00000000000000      Inscrição Estadual: 000000000
		Endereço: AV FAKE                                                 FAKE FAKE FAKE                        Telefone: 0000000000
																																							Informações Bancárias:
		Cidade:     FAKE                                                     UF: AM    CEP: 00000000                          Fax:
		Representante: FAKE                                                                                                    Telefone: 00 000000000       Banco :          Agência:     C.Corrente:
		Comprador:     FAKE DA SILVA                                                                                 Fax:
		Faturamento Embarque       Entrega                                          Condições de Pagamento : 30                                        Dias a partir da data de FATURAMENTO
		03/05/2023    03/05/2023   17/05/2023                                       Frete : CIF Transportadora :
	END TEXT

	The field that must be extracted from the text are:
		- Data
		- Número do pedido
		- CNPJ/CPF do fornecedor
		- CNPJ/CPF da empresa
		- Frete
		- Data de entrega

	The Schema and the output must be like this:
	{
		"purchase_order_number": "48740",
		"purchase_order_date": "03/05/2023",
		"cnpj_cpf_company": "00000000000000",
		"cnpj_cpf_supplier": "00000000000000",
		"freight_condition": "CIF",
		"delivery_date": "17/05/2023"
	}

	Do not include any fields that are not specified in the schema.
	Do not Generate any fields that are not specified in the schema.
	Do not generate any data if you don't have enough confidence in it.
`

const userPrompt = `
	TEXT
		{{.Text}}
	END TEXT`

// Header is the header of a purchase order
type Header struct {
	PurchaseOrderNumber string `json:"purchase_order_number"`
	PurchaseOrderDate   string `json:"purchase_order_date"`
	CnpjCpfCompany      string `json:"cnpj_cpf_company"`
	CnpjCpfSupplier     string `json:"cnpj_cpf_supplier"`
	FreightCondition    string `json:"freight_condition"`
	DeliveryDate        string `json:"delivery_date"`
}

func createUserPrompt(text string) (string, error) {
	type userInput struct {
		Text string
	}

	tmpl, err := template.New("parser").Parse(userPrompt)
	if err != nil {
		return "", fmt.Errorf("template parse error: %v", err)
	}

	headerParserInput := userInput{
		Text: text,
	}
	out := bytes.Buffer{}

	err = tmpl.Execute(&out, headerParserInput)
	if err != nil {
		return "", fmt.Errorf("template execute error: %v", err)
	}

	return out.String(), nil
}

// Config is the set of configuration parameters for the GPTClient
type Config struct {
	OpenAPIKey string
}

// Client is the client for the GPT API
type Client struct {
	client *openai.Client
}

// NewHeaderParser creates a new GPTClient
func NewHeaderParser(cfg Config) *Client {
	c := openai.NewClient(cfg.OpenAPIKey)
	return &Client{
		client: c,
	}
}

// HeaderParser parses the header of a purchase order
func (c *Client) HeaderParser(text string) (*Header, error) {
	userPrompt, err := createUserPrompt(text)
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
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userPrompt,
				},
			},
		})
	if err != nil {
		return nil, fmt.Errorf("gpt response error %v", err)

	}

	respContent := resp.Choices[0].Message.Content

	var header Header
	err = json.Unmarshal([]byte(respContent), &header)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal error %v", err)

	}
	return &header, nil
}
