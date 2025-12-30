package assets

import "embed"

//go:embed template.css propuesta.yaml html_template.yaml logo.svg presupuesto.yaml
var FS embed.FS

// GetCSS returns the embedded CSS template
func GetCSS() ([]byte, error) {
	return FS.ReadFile("template.css")
}

// GetPromptYAML returns the embedded prompt YAML
func GetPromptYAML() ([]byte, error) {
	return FS.ReadFile("propuesta.yaml")
}

// GetHTMLTemplateYAML returns the embedded HTML template YAML
func GetHTMLTemplateYAML() ([]byte, error) {
	return FS.ReadFile("html_template.yaml")
}

// GetLogo returns the embedded logo SVG
func GetLogo() ([]byte, error) {
	return FS.ReadFile("logo.svg")
}

// GetPresupuestoYAML returns the embedded presupuesto YAML
func GetPresupuestoYAML() ([]byte, error) {
	return FS.ReadFile("presupuesto.yaml")
}

