package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	DefaultModel = "claude-haiku-4-5-20251001"
)

var (
	ConfigDir  = filepath.Join(os.Getenv("HOME"), ".config", "orgmai")
	ConfigFile = filepath.Join(ConfigDir, "config.yaml")
)

type Config struct {
	ClaudeAPIKey string `yaml:"claude_api_key"`
	Model        string `yaml:"model"`
}

// Load carga la configuración desde el archivo YAML
func Load() (*Config, error) {
	// Crear directorio si no existe
	if err := os.MkdirAll(ConfigDir, 0755); err != nil {
		return nil, fmt.Errorf("error creando directorio de configuración: %w", err)
	}

	// Si el archivo no existe, retornar configuración por defecto
	if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
		return &Config{
			Model: DefaultModel,
		}, nil
	}

	// Leer archivo
	data, err := os.ReadFile(ConfigFile)
	if err != nil {
		return nil, fmt.Errorf("error leyendo archivo de configuración: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parseando configuración: %w", err)
	}

	// Establecer modelo por defecto si no está configurado
	if config.Model == "" {
		config.Model = DefaultModel
	}

	return &config, nil
}

// Save guarda la configuración en el archivo YAML
func Save(config *Config) error {
	// Crear directorio si no existe
	if err := os.MkdirAll(ConfigDir, 0755); err != nil {
		return fmt.Errorf("error creando directorio de configuración: %w", err)
	}

	// Serializar a YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error serializando configuración: %w", err)
	}

	// Escribir archivo
	if err := os.WriteFile(ConfigFile, data, 0644); err != nil {
		return fmt.Errorf("error escribiendo archivo de configuración: %w", err)
	}

	return nil
}

// GetAPIKey retorna la API key o error si no está configurada
func GetAPIKey() (string, error) {
	config, err := Load()
	if err != nil {
		return "", err
	}

	if config.ClaudeAPIKey == "" {
		return "", fmt.Errorf("API key no configurada. Ejecuta 'orgmai apikey'")
	}

	return config.ClaudeAPIKey, nil
}

// GetModel retorna el modelo configurado
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

// AvailableModels retorna la lista de modelos Claude disponibles
func AvailableModels() []string {
	return []string{
		"claude-haiku-4-5-20251001",
		"claude-sonnet-4-5-20250929",
		"claude-3-5-sonnet-20240620",
		"claude-3-opus-20240229",
	}
}



