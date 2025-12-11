package generator

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"orgmprop/assets"
	"orgmprop/internal/ai"
	"orgmprop/internal/config"
	"orgmprop/internal/logger"
)

// ProposalData represents the proposal data stored in JSON
type ProposalData struct {
	Titulo    string    `json:"titulo"`
	Subtitulo string    `json:"subtitulo"`
	Prompt    string    `json:"prompt"`
	Modelo    string    `json:"modelo"`
	Fecha     time.Time `json:"fecha"`
}

// GenerateProposal generates a complete proposal
func GenerateProposal(title, subtitle, prompt string, onProgress func(string)) (*ProposalData, string, error) {
	logger.Debug("Iniciando generación de propuesta: %s", title)

	// Get API key
	apiKey, err := config.GetAPIKey()
	if err != nil {
		return nil, "", fmt.Errorf("error obteniendo API key: %w", err)
	}

	// Get model
	model, err := config.GetModel()
	if err != nil {
		logger.Warn("Error obteniendo modelo, usando default: %v", err)
		model = config.DefaultModel
	}

	// Get prompt instructions
	promptInstructions, err := getPromptInstructions()
	if err != nil {
		return nil, "", fmt.Errorf("error obteniendo instrucciones de prompt: %w", err)
	}

	// Get HTML template instructions
	htmlInstructions, err := getHTMLInstructions()
	if err != nil {
		return nil, "", fmt.Errorf("error obteniendo instrucciones HTML: %w", err)
	}

	// Build system prompt
	systemPrompt := fmt.Sprintf(`%s

---

%s

---

IMPORTANTE: 
- Genera el HTML completo directamente desde el prompt del usuario.
- Usa el contenido generado según las reglas de prompt/propuesta.yaml.
- Formatea ese contenido en HTML usando la estructura de html/propuesta.yaml.
- NO incluyas CSS embebido ni etiquetas <style>. Usa un <link rel="stylesheet" href="template.css">.
- El HTML debe ser completo y listo para usar.
- NO uses markdown, solo HTML puro.
- NO incluyas bloques de código markdown.
- Asume que los assets (template.css, logo.svg/png) estarán en el mismo directorio que el HTML generado.`, promptInstructions, htmlInstructions)

	// Build user prompt
	userPrompt := fmt.Sprintf(`Título: %s
Subtítulo: %s

Prompt del usuario:
%s`, title, subtitle, prompt)

	// Create AI client
	client := ai.NewClient(apiKey, model)

	// Generate HTML
	logger.Debug("Generando HTML con IA...")
	var htmlContent string
	if onProgress != nil {
		htmlContent, err = client.GenerateProposalStream(context.Background(), systemPrompt, userPrompt, onProgress)
	} else {
		htmlContent, err = client.GenerateProposal(context.Background(), systemPrompt, userPrompt)
	}

	if err != nil {
		return nil, "", fmt.Errorf("error generando HTML: %w", err)
	}

	// Create proposal data
	proposalData := &ProposalData{
		Titulo:    title,
		Subtitulo: subtitle,
		Prompt:    prompt,
		Modelo:    model,
		Fecha:     time.Now(),
	}

	logger.Debug("Propuesta generada exitosamente")
	return proposalData, htmlContent, nil
}

// SaveProposal saves the proposal data and HTML to the current directory
func SaveProposal(data *ProposalData, htmlContent string) error {
	logger.Debug("Guardando propuesta en directorio actual")

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error obteniendo directorio actual: %w", err)
	}

	// Save JSON data
	jsonPath := filepath.Join(cwd, "propuesta.json")
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializando JSON: %w", err)
	}

	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		return fmt.Errorf("error guardando JSON: %w", err)
	}
	logger.Debug("JSON guardado en: %s", jsonPath)

	// Save HTML
	htmlPath := filepath.Join(cwd, "propuesta.html")
	if err := os.WriteFile(htmlPath, []byte(htmlContent), 0644); err != nil {
		return fmt.Errorf("error guardando HTML: %w", err)
	}
	logger.Debug("HTML guardado en: %s", htmlPath)

	// Copy logo
	if err := copyLogo(cwd); err != nil {
		logger.Warn("Error copiando logo: %v", err)
	}

	// Copy CSS
	if err := copyCSS(cwd); err != nil {
		logger.Warn("Error copiando CSS: %v", err)
	}

	return nil
}

// LoadProposal loads an existing proposal from the current directory
func LoadProposal() (*ProposalData, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo directorio actual: %w", err)
	}

	jsonPath := filepath.Join(cwd, "propuesta.json")
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo propuesta.json: %w", err)
	}

	var proposal ProposalData
	if err := json.Unmarshal(data, &proposal); err != nil {
		return nil, fmt.Errorf("error parseando propuesta.json: %w", err)
	}

	return &proposal, nil
}

// RegenerateProposal regenerates an existing proposal
func RegenerateProposal(data *ProposalData, onProgress func(string)) (string, error) {
	_, htmlContent, err := GenerateProposal(data.Titulo, data.Subtitulo, data.Prompt, onProgress)
	if err != nil {
		return "", err
	}

	return htmlContent, nil
}

// getPromptInstructions returns the prompt instructions from config or embedded assets
func getPromptInstructions() (string, error) {
	// First try to load from config directory
	configPath := config.GetConfigFilePath("propuesta.yaml")
	if data, err := os.ReadFile(configPath); err == nil {
		logger.Debug("Prompt YAML cargado desde config: %s", configPath)
		return string(data), nil
	}

	// Fall back to embedded assets
	data, err := assets.GetPromptYAML()
	if err != nil {
		return "", fmt.Errorf("error obteniendo prompt YAML embebido: %w", err)
	}

	logger.Debug("Prompt YAML cargado desde assets embebidos")
	return string(data), nil
}

// getHTMLInstructions returns the HTML template instructions from config or embedded assets
func getHTMLInstructions() (string, error) {
	// First try to load from config directory
	configPath := config.GetConfigFilePath("html_template.yaml")
	if data, err := os.ReadFile(configPath); err == nil {
		logger.Debug("HTML template YAML cargado desde config: %s", configPath)
		return string(data), nil
	}

	// Fall back to embedded assets
	data, err := assets.GetHTMLTemplateYAML()
	if err != nil {
		return "", fmt.Errorf("error obteniendo HTML template YAML embebido: %w", err)
	}

	logger.Debug("HTML template YAML cargado desde assets embebidos")
	return string(data), nil
}

// copyLogo copies the logo to the target directory
func copyLogo(targetDir string) error {
	// First try to load from config directory
	configPath := config.GetConfigFilePath("logo.svg")
	var logoData []byte
	var err error

	if data, readErr := os.ReadFile(configPath); readErr == nil {
		logoData = data
		logger.Debug("Logo cargado desde config: %s", configPath)
	} else {
		// Fall back to embedded assets
		logoData, err = assets.GetLogo()
		if err != nil {
			return fmt.Errorf("error obteniendo logo embebido: %w", err)
		}
		logger.Debug("Logo cargado desde assets embebidos")
	}

	// Determine extension based on config file
	logoExt := ".svg"
	pngPath := config.GetConfigFilePath("logo.png")
	if _, err := os.Stat(pngPath); err == nil {
		logoExt = ".png"
		logoData, _ = os.ReadFile(pngPath)
	}

	// Save logo to target directory
	targetPath := filepath.Join(targetDir, "logo"+logoExt)
	if err := os.WriteFile(targetPath, logoData, 0644); err != nil {
		return fmt.Errorf("error guardando logo: %w", err)
	}

	logger.Debug("Logo guardado en: %s", targetPath)
	return nil
}

// copyCSS copies the CSS to the target directory
func copyCSS(targetDir string) error {
	configPath := config.GetConfigFilePath("template.css")
	var cssData []byte
	var err error

	if data, readErr := os.ReadFile(configPath); readErr == nil {
		cssData = data
		logger.Debug("CSS cargado desde config: %s", configPath)
	} else {
		cssData, err = assets.GetCSS()
		if err != nil {
			return fmt.Errorf("error obteniendo CSS embebido: %w", err)
		}
		logger.Debug("CSS cargado desde assets embebidos")
	}

	targetPath := filepath.Join(targetDir, "template.css")
	if err := os.WriteFile(targetPath, cssData, 0644); err != nil {
		return fmt.Errorf("error guardando CSS: %w", err)
	}

	logger.Debug("CSS guardado en: %s", targetPath)
	return nil
}

