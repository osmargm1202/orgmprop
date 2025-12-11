package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Estilos
	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(false)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(false)

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Bold(false)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("204")).
			Bold(true)
)

// PrintInfo imprime un mensaje informativo
func PrintInfo(msg string) {
	fmt.Println(infoStyle.Render(msg))
}

// PrintSuccess imprime un mensaje de éxito
func PrintSuccess(msg string) {
	fmt.Println(successStyle.Render(msg))
}

// PrintWarning imprime un mensaje de advertencia
func PrintWarning(msg string) {
	fmt.Fprintf(os.Stderr, "%s\n", warningStyle.Render(msg))
}

// PrintError imprime un mensaje de error
func PrintError(msg string) {
	fmt.Fprintf(os.Stderr, "%s\n", errorStyle.Render(msg))
}

// Input solicita input del usuario
func Input(prompt, placeholder string) (string, error) {
	var value string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(prompt).
				Placeholder(placeholder).
				Value(&value),
		),
	)

	if err := form.Run(); err != nil {
		return "", err
	}

	return value, nil
}

// Select muestra opciones y permite seleccionar una
func Select(prompt string, options []string) (string, error) {
	var selected string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(prompt).
				Options(huh.NewOptions(options...)...).
				Value(&selected),
		),
	)

	if err := form.Run(); err != nil {
		return "", err
	}

	return selected, nil
}

// Confirm muestra una confirmación
func Confirm(message string) (bool, error) {
	var confirmed bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(message).
				Value(&confirmed),
		),
	)

	if err := form.Run(); err != nil {
		return false, err
	}

	return confirmed, nil
}



