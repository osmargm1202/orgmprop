package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// Custom theme for huh forms
func getTheme() *huh.Theme {
	t := huh.ThemeBase()

	t.Focused.Title = lipgloss.NewStyle().
		Foreground(ColorBlue).
		Bold(true)

	t.Focused.Description = lipgloss.NewStyle().
		Foreground(ColorSkyBlue)

	t.Focused.TextInput.Cursor = lipgloss.NewStyle().
		Foreground(ColorCyan)

	t.Focused.TextInput.Placeholder = lipgloss.NewStyle().
		Foreground(ColorGray)

	t.Focused.SelectSelector = lipgloss.NewStyle().
		Foreground(ColorCyan).
		SetString("> ")

	t.Focused.SelectedOption = lipgloss.NewStyle().
		Foreground(ColorBlue).
		Bold(true)

	t.Focused.UnselectedOption = lipgloss.NewStyle().
		Foreground(ColorSkyBlue)

	t.Blurred = t.Focused
	t.Blurred.Title = t.Blurred.Title.Foreground(ColorGray)

	return t
}

// PrintInfo prints an informational message
func PrintInfo(msg string) {
	fmt.Println(InfoStyle.Render(msg))
}

// PrintSuccess prints a success message
func PrintSuccess(msg string) {
	fmt.Println(SuccessStyle.Render("✓ " + msg))
}

// PrintWarning prints a warning message
func PrintWarning(msg string) {
	fmt.Fprintf(os.Stderr, "%s\n", WarningStyle.Render("⚠ "+msg))
}

// PrintError prints an error message
func PrintError(msg string) {
	fmt.Fprintf(os.Stderr, "%s\n", ErrorStyle.Render("✗ "+msg))
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
	).WithTheme(getTheme())

	if err := form.Run(); err != nil {
		return "", err
	}

	return value, nil
}

// InputPassword solicita input oculto del usuario
func InputPassword(prompt, placeholder string) (string, error) {
	var value string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(prompt).
				Placeholder(placeholder).
				EchoMode(huh.EchoModePassword).
				Value(&value),
		),
	).WithTheme(getTheme())

	if err := form.Run(); err != nil {
		return "", err
	}

	return value, nil
}

// TextArea solicita texto multilinea del usuario
func TextArea(prompt, placeholder string) (string, error) {
	var value string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title(prompt).
				Placeholder(placeholder).
				CharLimit(10000).
				Value(&value),
		),
	).WithTheme(getTheme())

	if err := form.Run(); err != nil {
		return "", err
	}

	return value, nil
}

// Select muestra opciones y permite seleccionar una
func Select(prompt string, options []string) (string, error) {
	var selected string

	opts := make([]huh.Option[string], len(options))
	for i, opt := range options {
		opts[i] = huh.NewOption(opt, opt)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(prompt).
				Options(opts...).
				Value(&selected),
		),
	).WithTheme(getTheme())

	if err := form.Run(); err != nil {
		return "", err
	}

	return selected, nil
}

// SelectWithKeys muestra opciones con valores diferentes a las keys
func SelectWithKeys(prompt string, options map[string]string) (string, error) {
	var selected string

	opts := make([]huh.Option[string], 0, len(options))
	for key, value := range options {
		opts = append(opts, huh.NewOption(key, value))
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(prompt).
				Options(opts...).
				Value(&selected),
		),
	).WithTheme(getTheme())

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
				Affirmative("Sí").
				Negative("No").
				Value(&confirmed),
		),
	).WithTheme(getTheme())

	if err := form.Run(); err != nil {
		return false, err
	}

	return confirmed, nil
}

// ProposalForm represents the form for creating a new proposal
type ProposalForm struct {
	Title    string
	Subtitle string
	Prompt   string
}

// NewProposalForm shows a form for creating a new proposal
func NewProposalForm() (*ProposalForm, error) {
	form := &ProposalForm{}

	f := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Título de la Propuesta").
				Placeholder("Ej: Diseño de Sistema Eléctrico").
				Value(&form.Title),
			huh.NewInput().
				Title("Subtítulo").
				Placeholder("Ej: Proyecto Residencial Plaza Norte").
				Value(&form.Subtitle),
		),
		huh.NewGroup(
			huh.NewText().
				Title("Descripción / Prompt").
				Placeholder("Describe el alcance, servicios, costos, tiempos...").
				CharLimit(10000).
				Value(&form.Prompt),
		),
	).WithTheme(getTheme())

	if err := f.Run(); err != nil {
		return nil, err
	}

	return form, nil
}

// ProjectForm represents the form for creating a new project
type ProjectForm struct {
	QuotationNumber string
	ProjectName     string
}

// NewProjectForm shows a form for creating a new project
func NewProjectForm() (*ProjectForm, error) {
	form := &ProjectForm{}

	f := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Número de Cotización").
				Placeholder("Ej: 001").
				Value(&form.QuotationNumber),
			huh.NewInput().
				Title("Nombre del Proyecto").
				Placeholder("Ej: Torre_Empresarial_Centro").
				Value(&form.ProjectName),
		),
	).WithTheme(getTheme())

	if err := f.Run(); err != nil {
		return nil, err
	}

	return form, nil
}

// NewPresupuestoForm shows a form for entering project description for budget generation
func NewPresupuestoForm() (string, error) {
	var descripcion string

	f := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Descripción del Proyecto").
				Placeholder("Describe el proyecto con todos los ítems, cantidades y precios...").
				CharLimit(10000).
				Value(&descripcion),
		),
	).WithTheme(getTheme())

	if err := f.Run(); err != nil {
		return "", err
	}

	return descripcion, nil
}

