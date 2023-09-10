package gpt

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/sashabaranov/go-openai"
)

// Config is the set of configuration parameters for the GPTClient
type Config struct {
	OpenAPIKey string
}

// Client is the client for the GPT API
type Client struct {
	client *openai.Client
}

func createUserPrompt(text string) (string, error) {
	type userInput struct {
		Text string
	}

	tmpl, err := template.New("parser").Parse(userHeaderPrompt)
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

// NewParser creates a new GPTClient
func NewParser(cfg Config) *Client {
	c := openai.NewClient(cfg.OpenAPIKey)
	return &Client{
		client: c,
	}
}
