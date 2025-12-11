package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	DefaultModel = "claude-sonnet-4-5-20250929"
	DotfilesURL  = "https://github.com/osmargm1202/dotfiles.git"
)

var (
	ConfigDir  = filepath.Join(os.Getenv("HOME"), ".config", "orgmprop")
	ConfigFile = filepath.Join(ConfigDir, "config.yaml")
)

// Config represents the application configuration
type Config struct {
	AnthropicAPIKey string `yaml:"anthropic_api_key"`
	Model           string `yaml:"model"`
	BaseFolder      string `yaml:"base_folder"`
}

// Load loads the configuration from YAML file
func Load() (*Config, error) {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(ConfigDir, 0755); err != nil {
		return nil, fmt.Errorf("error creando directorio de configuración: %w", err)
	}

	// If file doesn't exist, return default config
	if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
		return &Config{
			Model: DefaultModel,
		}, nil
	}

	// Read file
	data, err := os.ReadFile(ConfigFile)
	if err != nil {
		return nil, fmt.Errorf("error leyendo archivo de configuración: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parseando configuración: %w", err)
	}

	// Set default model if not configured
	if config.Model == "" {
		config.Model = DefaultModel
	}

	return &config, nil
}

// Save saves the configuration to YAML file
func Save(config *Config) error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(ConfigDir, 0755); err != nil {
		return fmt.Errorf("error creando directorio de configuración: %w", err)
	}

	// Serialize to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error serializando configuración: %w", err)
	}

	// Write file
	if err := os.WriteFile(ConfigFile, data, 0644); err != nil {
		return fmt.Errorf("error escribiendo archivo de configuración: %w", err)
	}

	return nil
}

// GetAPIKey returns the API key or error if not configured
func GetAPIKey() (string, error) {
	config, err := Load()
	if err != nil {
		return "", err
	}

	if config.AnthropicAPIKey == "" {
		return "", fmt.Errorf("API key no configurada. Ejecuta 'orgmprop config apikey'")
	}

	return config.AnthropicAPIKey, nil
}

// GetModel returns the configured model
func GetModel() (string, error) {
	config, err := Load()
	if err != nil {
		return DefaultModel, err
	}

	if config.Model == "" {
		return DefaultModel, nil
	}

	return config.Model, nil
}

// GetBaseFolder returns the configured base folder
func GetBaseFolder() (string, error) {
	config, err := Load()
	if err != nil {
		return "", err
	}

	if config.BaseFolder == "" {
		return "", fmt.Errorf("Carpeta base no configurada. Ejecuta 'orgmprop config folder'")
	}

	return config.BaseFolder, nil
}

// AvailableModels returns the list of available Anthropic models
func AvailableModels() []string {
	return []string{
		"claude-sonnet-4-5-20250929",
		"claude-haiku-4-5-20251001",
		"claude-3-5-sonnet-20241022",
		"claude-3-opus-20240229",
	}
}

// EnsureConfigFiles ensures all config files exist, downloading from dotfiles if needed
func EnsureConfigFiles() error {
	requiredFiles := []string{
		"template.css",
		"propuesta.yaml",
		"html_template.yaml",
		"logo.svg",
	}

	// Check if any file is missing
	missingFiles := []string{}
	for _, file := range requiredFiles {
		filePath := filepath.Join(ConfigDir, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			missingFiles = append(missingFiles, file)
		}
	}

	if len(missingFiles) == 0 {
		return nil
	}

	// Try to download from dotfiles
	if err := downloadDotfiles(); err != nil {
		// If dotfiles download fails, copy from embedded assets
		return copyFromEmbeddedAssets(missingFiles)
	}

	return nil
}

// downloadDotfiles downloads configuration from dotfiles repository
func downloadDotfiles() error {
	downloadsDir := filepath.Join(os.Getenv("HOME"), "Downloads")
	dotfilesDir := filepath.Join(downloadsDir, "dotfiles")

	// Create downloads dir if it doesn't exist
	if err := os.MkdirAll(downloadsDir, 0755); err != nil {
		return fmt.Errorf("error creando directorio Downloads: %w", err)
	}

	// Remove existing dotfiles dir
	os.RemoveAll(dotfilesDir)

	// Clone dotfiles repository
	cmd := exec.Command("git", "clone", "--depth", "1", DotfilesURL, dotfilesDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error clonando dotfiles: %w", err)
	}

	// Copy orgmprop config files
	orgmpropDir := filepath.Join(dotfilesDir, ".config", "orgmprop")
	if _, err := os.Stat(orgmpropDir); os.IsNotExist(err) {
		return fmt.Errorf("directorio orgmprop no encontrado en dotfiles")
	}

	// Copy files
	entries, err := os.ReadDir(orgmpropDir)
	if err != nil {
		return fmt.Errorf("error leyendo directorio orgmprop en dotfiles: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		src := filepath.Join(orgmpropDir, entry.Name())
		dst := filepath.Join(ConfigDir, entry.Name())

		// Don't overwrite existing config.yaml
		if entry.Name() == "config.yaml" {
			if _, err := os.Stat(dst); err == nil {
				continue
			}
		}

		data, err := os.ReadFile(src)
		if err != nil {
			continue
		}

		if err := os.WriteFile(dst, data, 0644); err != nil {
			continue
		}
	}

	return nil
}

// copyFromEmbeddedAssets copies missing files from embedded assets
func copyFromEmbeddedAssets(missingFiles []string) error {
	// This will be implemented using the assets package
	// For now, we'll return an error indicating manual setup is needed
	return fmt.Errorf("archivos de configuración faltantes: %v. Ejecuta el instalador o copia manualmente", missingFiles)
}

// GetConfigFilePath returns the path to a config file
func GetConfigFilePath(filename string) string {
	return filepath.Join(ConfigDir, filename)
}

// CopyTemplateFile copies a file to the config directory
func CopyTemplateFile(srcPath, destName string) error {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("error leyendo archivo: %w", err)
	}

	destPath := filepath.Join(ConfigDir, destName)
	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return fmt.Errorf("error escribiendo archivo: %w", err)
	}

	return nil
}

// ListMissingConfigFiles returns a list of missing configuration files
func ListMissingConfigFiles() []string {
	requiredFiles := []string{
		"template.css",
		"propuesta.yaml",
		"html_template.yaml",
		"logo.svg",
	}

	missing := []string{}
	for _, file := range requiredFiles {
		filePath := filepath.Join(ConfigDir, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			missing = append(missing, file)
		}
	}

	return missing
}

