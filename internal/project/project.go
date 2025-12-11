package project

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"orgmprop/internal/config"
	"orgmprop/internal/logger"
	"orgmprop/internal/ui"
)

// FolderStructure represents the folder.json structure
type FolderStructure struct {
	Tipos map[string]struct {
		Carpetas []string `json:"carpetas"`
	} `json:"tipos"`
}

// Default folder structure for projects
var defaultFolders = []string{
	"Comunicacion",
	"Diseño",
	"Estudios",
	"Calculos",
	"Levantamientos",
	"Entregas",
	"Recibido",
	"Oferta",
}

// CreateProject creates a new project structure
func CreateProject(quotationNumber, projectName string) (string, error) {
	logger.Debug("Creando proyecto: %s - %s", quotationNumber, projectName)

	// Get base folder
	baseFolder, err := config.GetBaseFolder()
	if err != nil {
		return "", fmt.Errorf("carpeta base no configurada: %w", err)
	}

	// Sanitize project name
	projectName = sanitizeName(projectName)
	folderName := fmt.Sprintf("%s-%s", quotationNumber, projectName)

	// Create project directory
	projectPath := filepath.Join(baseFolder, folderName)
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return "", fmt.Errorf("error creando directorio del proyecto: %w", err)
	}

	logger.Debug("Directorio del proyecto creado: %s", projectPath)

	// Create subfolders
	folders := getFolderStructure()
	for _, folder := range folders {
		folderPath := filepath.Join(projectPath, folder)
		if err := os.MkdirAll(folderPath, 0755); err != nil {
			logger.Warn("Error creando carpeta %s: %v", folder, err)
		} else {
			logger.Debug("Carpeta creada: %s", folderPath)
		}
	}

	// Return path to Oferta folder
	ofertaPath := filepath.Join(projectPath, "Oferta")
	logger.Debug("Proyecto creado exitosamente, carpeta Oferta: %s", ofertaPath)

	return ofertaPath, nil
}

// ListProjects lists all projects in the base folder
func ListProjects() ([]string, error) {
	logger.Debug("Listando proyectos")

	// Get base folder
	baseFolder, err := config.GetBaseFolder()
	if err != nil {
		return nil, fmt.Errorf("carpeta base no configurada: %w", err)
	}

	// Read directory
	entries, err := os.ReadDir(baseFolder)
	if err != nil {
		return nil, fmt.Errorf("error leyendo directorio base: %w", err)
	}

	var projects []string
	for _, entry := range entries {
		if entry.IsDir() {
			projects = append(projects, entry.Name())
		}
	}

	// Sort by name
	sort.Strings(projects)

	logger.Debug("Encontrados %d proyectos", len(projects))
	return projects, nil
}

// GetProjectOfertaPath returns the path to the Oferta folder for a project
func GetProjectOfertaPath(projectName string) (string, error) {
	// Get base folder
	baseFolder, err := config.GetBaseFolder()
	if err != nil {
		return "", fmt.Errorf("carpeta base no configurada: %w", err)
	}

	projectPath := filepath.Join(baseFolder, projectName)
	ofertaPath := filepath.Join(projectPath, "Oferta")

	// Create Oferta folder if it doesn't exist
	if err := os.MkdirAll(ofertaPath, 0755); err != nil {
		return "", fmt.Errorf("error creando carpeta Oferta: %w", err)
	}

	return ofertaPath, nil
}

// EnsureProjectStructure ensures all folders exist for a project
func EnsureProjectStructure(projectName string) error {
	logger.Debug("Verificando estructura del proyecto: %s", projectName)

	// Get base folder
	baseFolder, err := config.GetBaseFolder()
	if err != nil {
		return fmt.Errorf("carpeta base no configurada: %w", err)
	}

	projectPath := filepath.Join(baseFolder, projectName)

	// Create project directory if it doesn't exist
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf("error creando directorio del proyecto: %w", err)
	}

	// Create subfolders
	folders := getFolderStructure()
	for _, folder := range folders {
		folderPath := filepath.Join(projectPath, folder)
		if err := os.MkdirAll(folderPath, 0755); err != nil {
			logger.Warn("Error creando carpeta %s: %v", folder, err)
		}
	}

	logger.Debug("Estructura del proyecto verificada")
	return nil
}

// ProposalSummary represents a summary of a proposal
type ProposalSummary struct {
	Project   string
	Title     string
	Subtitle  string
	Date      time.Time
	FilePath  string
	HTMLPath  string
}

// GetProposalSummaries scans all projects for proposals
func GetProposalSummaries() ([]ui.ProposalSummary, error) {
	logger.Debug("Escaneando propuestas")

	// Get base folder
	baseFolder, err := config.GetBaseFolder()
	if err != nil {
		return nil, fmt.Errorf("carpeta base no configurada: %w", err)
	}

	var summaries []ui.ProposalSummary

	// Read all project directories
	entries, err := os.ReadDir(baseFolder)
	if err != nil {
		return nil, fmt.Errorf("error leyendo directorio base: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		projectName := entry.Name()
		ofertaPath := filepath.Join(baseFolder, projectName, "Oferta")

		// Check if Oferta folder exists
		if _, err := os.Stat(ofertaPath); os.IsNotExist(err) {
			continue
		}

		// Check for propuesta.json
		jsonPath := filepath.Join(ofertaPath, "propuesta.json")
		if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
			continue
		}

		// Read proposal data
		data, err := os.ReadFile(jsonPath)
		if err != nil {
			logger.Warn("Error leyendo %s: %v", jsonPath, err)
			continue
		}

		var proposal struct {
			Titulo    string    `json:"titulo"`
			Subtitulo string    `json:"subtitulo"`
			Fecha     time.Time `json:"fecha"`
		}

		if err := json.Unmarshal(data, &proposal); err != nil {
			logger.Warn("Error parseando %s: %v", jsonPath, err)
			continue
		}

		htmlPath := filepath.Join(ofertaPath, "propuesta.html")

		summaries = append(summaries, ui.ProposalSummary{
			Project:  projectName,
			Title:    proposal.Titulo,
			Date:     proposal.Fecha.Format("2006-01-02"),
			FilePath: htmlPath,
		})
	}

	// Sort by date (newest first)
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Date > summaries[j].Date
	})

	logger.Debug("Encontradas %d propuestas", len(summaries))
	return summaries, nil
}

// getFolderStructure returns the folder structure to use
func getFolderStructure() []string {
	// Try to load from folder.json in config directory
	configPath := config.GetConfigFilePath("folder.json")
	if data, err := os.ReadFile(configPath); err == nil {
		var structure FolderStructure
		if err := json.Unmarshal(data, &structure); err == nil {
			if proyectos, ok := structure.Tipos["Proyectos"]; ok {
				logger.Debug("Estructura de carpetas cargada desde config")
				return proyectos.Carpetas
			}
		}
	}

	logger.Debug("Usando estructura de carpetas por defecto")
	return defaultFolders
}

// sanitizeName sanitizes a name for use in file paths
func sanitizeName(name string) string {
	// Replace spaces with underscores
	name = strings.ReplaceAll(name, " ", "_")
	
	// Remove or replace invalid characters
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		name = strings.ReplaceAll(name, char, "_")
	}

	return name
}

// ChangeToOfertaDirectory changes to the Oferta directory of a project
func ChangeToOfertaDirectory(projectName string) error {
	ofertaPath, err := GetProjectOfertaPath(projectName)
	if err != nil {
		return err
	}

	if err := os.Chdir(ofertaPath); err != nil {
		return fmt.Errorf("error cambiando a directorio Oferta: %w", err)
	}

	logger.Debug("Cambiado a directorio: %s", ofertaPath)
	return nil
}

