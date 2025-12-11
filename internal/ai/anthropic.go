package ai

import (
	"context"
	"fmt"
	"strings"

	"orgmprop/internal/logger"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// Client encapsulates the Anthropic client
type Client struct {
	client *anthropic.Client
	model  string
}

// NewClient creates a new AI client
func NewClient(apiKey, model string) *Client {
	logger.Debug("Creando cliente Anthropic con modelo: %s", model)
	client := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &Client{
		client: client,
		model:  model,
	}
}

// GenerateProposal generates a proposal HTML using the AI model
func (c *Client) GenerateProposal(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	logger.Debug("Generando propuesta con modelo: %s", c.model)
	logger.Debug("System prompt length: %d", len(systemPrompt))
	logger.Debug("User prompt length: %d", len(userPrompt))

	// Create message params using the F helper
	params := anthropic.MessageNewParams{
		Model:     anthropic.F(anthropic.Model(c.model)),
		MaxTokens: anthropic.F(int64(8192)),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(userPrompt)),
		}),
		System: anthropic.F([]anthropic.TextBlockParam{
			{
				Type: anthropic.F(anthropic.TextBlockParamTypeText),
				Text: anthropic.F(systemPrompt),
			},
		}),
	}

	// Send request
	logger.Debug("Enviando solicitud a Anthropic...")
	resp, err := c.client.Messages.New(ctx, params)
	if err != nil {
		logger.Error("Error en solicitud a Anthropic: %v", err)
		return "", fmt.Errorf("error en la solicitud a Anthropic: %w", err)
	}

	logger.Debug("Respuesta recibida, procesando contenido...")

	// Extract content from response
	if len(resp.Content) == 0 {
		logger.Error("Respuesta vacía de Anthropic")
		return "", fmt.Errorf("respuesta vacía de Anthropic")
	}

	// The content can be text or blocks, extract text
	var content strings.Builder
	for _, block := range resp.Content {
		if block.Type == anthropic.ContentBlockTypeText {
			content.WriteString(block.Text)
		}
	}

	result := content.String()
	logger.Debug("Contenido generado, longitud: %d", len(result))

	// Clean up HTML if needed
	result = cleanHTMLResponse(result)

	return result, nil
}

// GenerateProposalStream generates a proposal HTML using streaming
func (c *Client) GenerateProposalStream(ctx context.Context, systemPrompt, userPrompt string, onChunk func(string)) (string, error) {
	logger.Debug("Generando propuesta con streaming, modelo: %s", c.model)

	// Create message params using the F helper
	params := anthropic.MessageNewParams{
		Model:     anthropic.F(anthropic.Model(c.model)),
		MaxTokens: anthropic.F(int64(8192)),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(userPrompt)),
		}),
		System: anthropic.F([]anthropic.TextBlockParam{
			{
				Type: anthropic.F(anthropic.TextBlockParamTypeText),
				Text: anthropic.F(systemPrompt),
			},
		}),
	}

	// Create stream
	stream := c.client.Messages.NewStreaming(ctx, params)

	// Read and accumulate response
	var fullResponse strings.Builder

	for stream.Next() {
		event := stream.Current()

		// Process content_block_delta events
		if event.Type == anthropic.MessageStreamEventTypeContentBlockDelta {
			// Type assert Delta to access its fields
			if delta, ok := event.Delta.(anthropic.ContentBlockDeltaEventDelta); ok {
				if delta.Type == anthropic.ContentBlockDeltaEventDeltaTypeTextDelta {
					fullResponse.WriteString(delta.Text)
					if onChunk != nil {
						onChunk(delta.Text)
					}
				}
			}
		}
	}

	if err := stream.Err(); err != nil {
		logger.Error("Error leyendo stream: %v", err)
		return "", fmt.Errorf("error leyendo stream: %w", err)
	}

	result := fullResponse.String()
	logger.Debug("Streaming completado, longitud: %d", len(result))

	// Clean up HTML
	result = cleanHTMLResponse(result)

	return result, nil
}

// cleanHTMLResponse cleans up the HTML response from the AI
func cleanHTMLResponse(html string) string {
	// Remove markdown code fences if present
	html = strings.TrimSpace(html)
	if strings.HasPrefix(html, "```html") {
		html = html[7:]
	}
	if strings.HasPrefix(html, "```") {
		html = html[3:]
	}
	if strings.HasSuffix(html, "```") {
		html = html[:len(html)-3]
	}
	html = strings.TrimSpace(html)

	return html
}
