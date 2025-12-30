package generator

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"orgmprop/assets"
	"orgmprop/internal/ai"
	"orgmprop/internal/config"
	"orgmprop/internal/logger"

	"gopkg.in/yaml.v3"
)

// PresupuestoYAML represents the structure of presupuesto.yaml
type PresupuestoYAML struct {
	System       string `yaml:"system"`
	UserTemplate string `yaml:"user_template"`
	Variables    []struct {
		Name        string `yaml:"name"`
		Type        string `yaml:"type"`
		Required    bool   `yaml:"required"`
		Description string `yaml:"description"`
	} `yaml:"variables"`
	OutputFormat struct {
		Type   string `yaml:"type"`
		Schema struct {
			Type       string   `yaml:"type"`
			Required   []string `yaml:"required"`
			Properties map[string]interface{} `yaml:"properties"`
		} `yaml:"schema"`
	} `yaml:"output_format"`
	NotasTecnicas string `yaml:"notas_tecnicas"`
}

// GeneratePresupuesto generates a budget JSON using the presupuesto.yaml prompt
func GeneratePresupuesto(descripcionProyecto string, onProgress func(string)) ([]byte, error) {
	logger.Debug("Iniciando generación de presupuesto")

	// Get API key
	apiKey, err := config.GetAPIKey()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo API key: %w", err)
	}

	// Get model
	model, err := config.GetModel()
	if err != nil {
		logger.Warn("Error obteniendo modelo, usando default: %v", err)
		model = config.DefaultModel
	}

	// Load presupuesto YAML
	yamlData, err := getPresupuestoYAML()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo presupuesto YAML: %w", err)
	}

	// Parse YAML and extract JSON example
	systemPrompt, userTemplate, ejemploJSON, err := parsePresupuestoYAML(yamlData)
	if err != nil {
		return nil, fmt.Errorf("error parseando presupuesto YAML: %w", err)
	}

	// Build user prompt by replacing variables
	userPrompt := strings.ReplaceAll(userTemplate, "{descripcion_proyecto}", descripcionProyecto)
	userPrompt = strings.ReplaceAll(userPrompt, "{ejemplo_json}", ejemploJSON)

	logger.Debug("System prompt length: %d", len(systemPrompt))
	logger.Debug("User prompt length: %d", len(userPrompt))

	// Create AI client
	client := ai.NewClient(apiKey, model)

	// Generate JSON
	logger.Debug("Generando JSON con IA...")
	var jsonContent string
	if onProgress != nil {
		jsonContent, err = client.GenerateProposalStream(context.Background(), systemPrompt, userPrompt, onProgress)
	} else {
		jsonContent, err = client.GenerateProposal(context.Background(), systemPrompt, userPrompt)
	}

	if err != nil {
		return nil, fmt.Errorf("error generando JSON: %w", err)
	}

	// Clean and validate JSON
	logger.Debug("Respuesta de IA antes de limpiar (primeros 200 chars): %s", jsonContent[:min(200, len(jsonContent))])
	jsonContent = cleanJSONResponse(jsonContent)
	logger.Debug("JSON después de limpiar (primeros 200 chars): %s", jsonContent[:min(200, len(jsonContent))])
	
	// Validate JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonContent), &jsonData); err != nil {
		logger.Error("JSON inválido generado. Error: %v", err)
		logger.Error("Contenido recibido (primeros 500 chars): %s", jsonContent[:min(500, len(jsonContent))])
		return nil, fmt.Errorf("error validando JSON generado: %w", err)
	}

	// Format JSON with indentation
	formattedJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error formateando JSON: %w", err)
	}

	logger.Debug("Presupuesto generado exitosamente")
	return formattedJSON, nil
}

// PresupuestoPromptData represents the prompt data stored in a text file
type PresupuestoPromptData struct {
	Prompt    string
	Timestamp time.Time
}

// SavePresupuestoPrompt saves the prompt to a file in the current directory
func SavePresupuestoPrompt(prompt string) error {
	logger.Debug("Guardando prompt de presupuesto")

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error obteniendo directorio actual: %w", err)
	}

	// Check if prompt file exists
	promptPath := filepath.Join(cwd, "presupuesto_prompt.txt")
	existingPrompt := ""
	if data, err := os.ReadFile(promptPath); err == nil {
		existingPrompt = string(data)
		logger.Debug("Prompt existente encontrado, agregando nuevo contenido")
	}

	// Combine prompts
	var finalPrompt string
	if existingPrompt != "" {
		finalPrompt = existingPrompt + "\n\n--- NUEVA SOLICITUD ---\n\n" + prompt
	} else {
		finalPrompt = prompt
	}

	// Save prompt
	if err := os.WriteFile(promptPath, []byte(finalPrompt), 0644); err != nil {
		return fmt.Errorf("error guardando prompt: %w", err)
	}

	logger.Debug("Prompt guardado en: %s", promptPath)
	return nil
}

// LoadPresupuestoPrompt loads the existing prompt from the current directory
func LoadPresupuestoPrompt() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error obteniendo directorio actual: %w", err)
	}

	promptPath := filepath.Join(cwd, "presupuesto_prompt.txt")
	data, err := os.ReadFile(promptPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // No existe, retornar vacío
		}
		return "", fmt.Errorf("error leyendo prompt: %w", err)
	}

	return string(data), nil
}

// SavePresupuesto saves the budget JSON to the current directory
func SavePresupuesto(jsonData []byte) error {
	logger.Debug("Guardando presupuesto en directorio actual")

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error obteniendo directorio actual: %w", err)
	}

	// Save JSON
	jsonPath := filepath.Join(cwd, "presupuesto.json")
	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		return fmt.Errorf("error guardando JSON: %w", err)
	}

	logger.Debug("Presupuesto guardado en: %s", jsonPath)
	return nil
}

// getPresupuestoYAML returns the presupuesto YAML from config or embedded assets
func getPresupuestoYAML() ([]byte, error) {
	// First try to load from config directory
	configPath := config.GetPresupuestoYAMLFilePath()
	if data, err := os.ReadFile(configPath); err == nil {
		logger.Debug("Presupuesto YAML cargado desde config: %s", configPath)
		return data, nil
	}

	// Fall back to embedded assets
	data, err := assets.GetPresupuestoYAML()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo presupuesto YAML embebido: %w", err)
	}

	logger.Debug("Presupuesto YAML cargado desde assets embebidos")
	return data, nil
}

// parsePresupuestoYAML parses the YAML file and extracts system prompt, user template, and JSON example
func parsePresupuestoYAML(data []byte) (systemPrompt, userTemplate, ejemploJSON string, err error) {
	// Convert to string for processing
	content := string(data)

	// Find the marker for the JSON example
	marker := "aqui debajo dejo la cotizaicon de ejmeplo:"
	markerIndex := strings.Index(content, marker)
	if markerIndex == -1 {
		return "", "", "", fmt.Errorf("marcador de ejemplo JSON no encontrado")
	}

	// Split content: YAML part and JSON part
	yamlPart := content[:markerIndex]
	jsonPart := content[markerIndex+len(marker):]

	// Parse YAML part
	var yamlData PresupuestoYAML
	if err := yaml.Unmarshal([]byte(yamlPart), &yamlData); err != nil {
		return "", "", "", fmt.Errorf("error parseando YAML: %w", err)
	}

	// Extract JSON example (clean it)
	ejemploJSON = strings.TrimSpace(jsonPart)
	// Remove leading/trailing whitespace and newlines
	ejemploJSON = strings.Trim(ejemploJSON, "\n\r\t ")

	logger.Debug("System prompt extraído, longitud: %d", len(yamlData.System))
	logger.Debug("User template extraído, longitud: %d", len(yamlData.UserTemplate))
	logger.Debug("Ejemplo JSON extraído, longitud: %d", len(ejemploJSON))

	return yamlData.System, yamlData.UserTemplate, ejemploJSON, nil
}

// cleanJSONResponse cleans up the JSON response from the AI
func cleanJSONResponse(response string) string {
	// Remove markdown code fences if present
	response = strings.TrimSpace(response)
	
	// Remove ```json or ``` at the beginning
	if strings.HasPrefix(response, "```json") {
		response = strings.TrimPrefix(response, "```json")
	} else if strings.HasPrefix(response, "```") {
		response = strings.TrimPrefix(response, "```")
	}
	
	// Remove ``` at the end
	if strings.HasSuffix(response, "```") {
		response = strings.TrimSuffix(response, "```")
	}
	
	response = strings.TrimSpace(response)
	
	// Find the first '{' or '[' which should be the start of JSON
	startIndex := -1
	for i, char := range response {
		if char == '{' || char == '[' {
			startIndex = i
			break
		}
	}
	
	if startIndex == -1 {
		// No JSON found, return as is (will fail validation)
		return response
	}
	
	// Find the last '}' or ']' which should be the end of JSON
	endIndex := -1
	for i := len(response) - 1; i >= 0; i-- {
		char := response[i]
		if char == '}' || char == ']' {
			endIndex = i + 1
			break
		}
	}
	
	if endIndex == -1 || endIndex <= startIndex {
		// No valid end found, return from start to end
		return response[startIndex:]
	}
	
	// Extract only the JSON part
	json := response[startIndex:endIndex]
	return strings.TrimSpace(json)
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

