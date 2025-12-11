package chat

import (
	"context"
	"fmt"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/charmbracelet/lipgloss"
	"orgmai/internal/conversation"
)

const (
	SystemPrompt = "Responde de forma compacta sin espacios innecesarios entre líneas. Usa formato markdown optimizado para terminal Linux."
)

var (
	assistantStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39"))
)

// Client encapsula el cliente de Anthropic
type Client struct {
	client anthropic.Client
	model  string
}

// NewClient crea un nuevo cliente de chat
func NewClient(apiKey, model string) *Client {
	return &Client{
		client: anthropic.NewClient(option.WithAPIKey(apiKey)),
		model:  model,
	}
}

// Send envía mensajes a Claude y retorna la respuesta
func (c *Client) Send(ctx context.Context, messages []conversation.Message) (string, error) {
	// Convertir mensajes al formato de Anthropic
	anthropicMessages := make([]anthropic.MessageParam, 0, len(messages))
	
	for _, msg := range messages {
		if msg.Role == "user" {
			anthropicMessages = append(anthropicMessages, anthropic.NewUserMessage(
				anthropic.NewTextBlock(msg.Content),
			))
		} else if msg.Role == "assistant" {
			anthropicMessages = append(anthropicMessages, anthropic.NewAssistantMessage(
				anthropic.NewTextBlock(msg.Content),
			))
		}
	}

	// Crear system prompt como TextBlockParam
	systemBlock := anthropic.TextBlockParam{Text: SystemPrompt}

	// Crear parámetros de mensaje
	params := anthropic.MessageNewParams{
		Model:     anthropic.Model(c.model),
		MaxTokens: 4096,
		Messages:  anthropicMessages,
		System:    []anthropic.TextBlockParam{systemBlock},
	}

	// Enviar solicitud
	resp, err := c.client.Messages.New(ctx, params)
	if err != nil {
		return "", fmt.Errorf("error en la solicitud a Claude: %w", err)
	}

	// Extraer contenido de la respuesta
	if len(resp.Content) == 0 {
		return "", fmt.Errorf("respuesta vacía de Claude")
	}

	// El contenido puede ser texto o bloques, extraer texto
	var content strings.Builder
	for _, block := range resp.Content {
		if block.Type == "text" {
			textBlock := block.AsText()
			content.WriteString(textBlock.Text)
		}
	}

	return content.String(), nil
}

// SendStream envía mensajes a Claude y muestra la respuesta en streaming
func (c *Client) SendStream(ctx context.Context, messages []conversation.Message) (string, error) {
	// Convertir mensajes al formato de Anthropic
	anthropicMessages := make([]anthropic.MessageParam, 0, len(messages))
	
	for _, msg := range messages {
		if msg.Role == "user" {
			anthropicMessages = append(anthropicMessages, anthropic.NewUserMessage(
				anthropic.NewTextBlock(msg.Content),
			))
		} else if msg.Role == "assistant" {
			anthropicMessages = append(anthropicMessages, anthropic.NewAssistantMessage(
				anthropic.NewTextBlock(msg.Content),
			))
		}
	}

	// Crear system prompt como TextBlockParam
	systemBlock := anthropic.TextBlockParam{Text: SystemPrompt}

	// Crear parámetros de mensaje
	params := anthropic.MessageNewParams{
		Model:     anthropic.Model(c.model),
		MaxTokens: 4096,
		Messages:  anthropicMessages,
		System:    []anthropic.TextBlockParam{systemBlock},
	}

	// Crear stream
	stream := c.client.Messages.NewStreaming(ctx, params)

	// Leer y mostrar respuesta en streaming
	var fullResponse strings.Builder
	fmt.Print("\n")
	
	for stream.Next() {
		event := stream.Current()
		
		// Procesar diferentes tipos de eventos según el tipo
		if event.Type == "content_block_delta" {
			if event.Delta.Type == "text_delta" {
				text := event.Delta.Text
				fmt.Print(text)
				fullResponse.WriteString(text)
			}
		}
	}

	if err := stream.Err(); err != nil {
		return "", fmt.Errorf("error leyendo stream: %w", err)
	}

	fmt.Print("\n\n")
	return fullResponse.String(), nil
}

// RenderMarkdown renderiza markdown simple en la terminal
func RenderMarkdown(text string) {
	// Por ahora simplemente imprimir el texto
	// En el futuro se podría usar una librería de markdown rendering
	fmt.Println(text)
}

