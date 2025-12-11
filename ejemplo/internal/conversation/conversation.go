package conversation

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	MaxConversations = 10
)

var (
	ConversationsDir = filepath.Join(os.Getenv("HOME"), ".config", "orgmai", "conversations")
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// sanitizeFilename sanitiza el título para usarlo como nombre de archivo
func sanitizeFilename(title string) string {
	// Remover caracteres especiales
	re := regexp.MustCompile(`[^\w\s-]`)
	title = re.ReplaceAllString(title, "")
	
	// Reemplazar espacios múltiples y guiones con un solo guion
	re = regexp.MustCompile(`[-\s]+`)
	title = re.ReplaceAllString(title, "-")
	
	// Limitar longitud
	if len(title) > 50 {
		title = title[:50]
	}
	
	return title
}

// areWordsSimilar verifica si dos palabras son similares
func areWordsSimilar(word1, word2 string) bool {
	word1Lower := strings.ToLower(word1)
	word2Lower := strings.ToLower(word2)

	if word1Lower == word2Lower {
		return true
	}

	if len(word1Lower) >= 4 && len(word2Lower) >= 4 {
		if strings.HasPrefix(word1Lower, word2Lower[:4]) || strings.HasPrefix(word2Lower, word1Lower[:4]) {
			return true
		}
	}

	return false
}

// GenerateTitle genera un título tomando las 3 palabras más grandes del primer prompt
func GenerateTitle(messages []Message) string {
	if len(messages) == 0 {
		return "conversacion"
	}

	// Buscar primer mensaje del usuario
	var firstUserMessage string
	for _, msg := range messages {
		if msg.Role == "user" {
			firstUserMessage = msg.Content
			break
		}
	}

	if firstUserMessage == "" {
		return "conversacion"
	}

	// Extraer palabras significativas
	words := strings.Fields(firstUserMessage)
	var meaningfulWords []string
	for _, w := range words {
		cleaned := strings.Trim(w, ".,!?;:()[]{}\"'")
		if len(cleaned) >= 3 {
			meaningfulWords = append(meaningfulWords, cleaned)
		}
	}

	// Ordenar por longitud (más grandes primero)
	sort.Slice(meaningfulWords, func(i, j int) bool {
		return len(meaningfulWords[i]) > len(meaningfulWords[j])
	})

	// Seleccionar palabras únicas (no repetidas ni similares)
	var titleWords []string
	for _, word := range meaningfulWords {
		isDuplicate := false
		for _, existing := range titleWords {
			if areWordsSimilar(word, existing) {
				isDuplicate = true
				break
			}
		}
		if !isDuplicate {
			titleWords = append(titleWords, word)
			if len(titleWords) >= 3 {
				break
			}
		}
	}

	// Si no hay suficientes palabras únicas, usar las primeras palabras del mensaje
	if len(titleWords) < 3 {
		allWords := strings.Fields(firstUserMessage)
		for _, word := range allWords {
			cleaned := strings.Trim(word, ".,!?;:()[]{}\"'")
			if len(cleaned) >= 3 {
				isDuplicate := false
				for _, existing := range titleWords {
					if areWordsSimilar(cleaned, existing) {
						isDuplicate = true
						break
					}
				}
				if !isDuplicate {
					titleWords = append(titleWords, cleaned)
					if len(titleWords) >= 3 {
						break
					}
				}
			}
		}
	}

	title := strings.Join(titleWords, " ")
	if title == "" {
		title = "conversacion"
	}

	return sanitizeFilename(title)
}

// Save guarda la conversación en un archivo MD
func Save(messages []Message, title, timestamp string) error {
	// Crear directorio si no existe
	if err := os.MkdirAll(ConversationsDir, 0755); err != nil {
		return fmt.Errorf("error creando directorio de conversaciones: %w", err)
	}

	filename := fmt.Sprintf("%s-%s.md", timestamp, title)
	filepath := filepath.Join(ConversationsDir, filename)

	// Parsear timestamp
	dt, err := time.Parse("2006-01-02-15-04-05", timestamp)
	dateStr := timestamp
	if err == nil {
		dateStr = dt.Format("2006-01-02 15:04:05")
	}

	// Escribir archivo
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("error creando archivo: %w", err)
	}
	defer file.Close()

	// Escribir encabezado
	fmt.Fprintf(file, "# %s\n\n", title)
	fmt.Fprintf(file, "**Fecha:** %s\n\n", dateStr)
	fmt.Fprintf(file, "---\n\n")

	// Escribir mensajes
	for _, msg := range messages {
		if msg.Role == "user" {
			fmt.Fprintf(file, "## Usuario\n\n%s\n\n", msg.Content)
		} else if msg.Role == "assistant" {
			fmt.Fprintf(file, "## Asistente\n\n%s\n\n", msg.Content)
		}
		fmt.Fprintf(file, "---\n\n")
	}

	return nil
}

// Load carga una conversación desde un archivo MD
func Load(filename string) ([]Message, error) {
	filepath := filepath.Join(ConversationsDir, filename)

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, fmt.Errorf("archivo no encontrado: %s", filename)
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo archivo: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	var messages []Message
	var currentRole string
	var currentContent []string
	inContent := false

	for _, line := range lines {
		if strings.HasPrefix(line, "## Usuario") {
			if currentRole != "" && len(currentContent) > 0 {
				messages = append(messages, Message{
					Role:    currentRole,
					Content: strings.Join(currentContent, "\n"),
				})
			}
			currentRole = "user"
			currentContent = []string{}
			inContent = true
		} else if strings.HasPrefix(line, "## Asistente") {
			if currentRole != "" && len(currentContent) > 0 {
				messages = append(messages, Message{
					Role:    currentRole,
					Content: strings.Join(currentContent, "\n"),
				})
			}
			currentRole = "assistant"
			currentContent = []string{}
			inContent = true
		} else if strings.HasPrefix(line, "---") {
			continue
		} else if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "**") {
			continue
		} else if inContent {
			currentContent = append(currentContent, line)
		}
	}

	// Agregar último mensaje
	if currentRole != "" && len(currentContent) > 0 {
		messages = append(messages, Message{
			Role:    currentRole,
			Content: strings.TrimSpace(strings.Join(currentContent, "\n")),
		})
	}

	return messages, nil
}

// List lista las últimas conversaciones (máximo 10)
func List() ([]string, error) {
	if _, err := os.Stat(ConversationsDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	entries, err := os.ReadDir(ConversationsDir)
	if err != nil {
		return nil, fmt.Errorf("error leyendo directorio: %w", err)
	}

	var files []struct {
		name  string
		mtime time.Time
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		files = append(files, struct {
			name  string
			mtime time.Time
		}{
			name:  entry.Name(),
			mtime: info.ModTime(),
		})
	}

	// Ordenar por fecha de modificación (más reciente primero)
	sort.Slice(files, func(i, j int) bool {
		return files[i].mtime.After(files[j].mtime)
	})

	// Limitar a MaxConversations
	var result []string
	for i, file := range files {
		if i >= MaxConversations {
			break
		}
		result = append(result, file.name)
	}

	return result, nil
}



