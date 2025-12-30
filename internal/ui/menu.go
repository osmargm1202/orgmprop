package ui

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

// MenuOption represents a menu option
type MenuOption struct {
	Label string
	Value string
}

// MainMenuOptions returns the main menu options
func MainMenuOptions() []MenuOption {
	return []MenuOption{
		{Label: "📝 Nueva Propuesta", Value: "new"},
		{Label: "💰 Generar Presupuesto", Value: "presupuesto"},
		{Label: "📂 Crear Proyecto", Value: "proyecto"},
		{Label: "📋 Listar Proyectos", Value: "list"},
		{Label: "📊 Resumen de Propuestas", Value: "resumen"},
		{Label: "💰 Resumen de Presupuestos", Value: "resumen_presupuestos"},
		{Label: "⚙️  Configuración", Value: "config"},
		{Label: "❌ Salir", Value: "exit"},
	}
}

// ConfigMenuOptions returns the config menu options
func ConfigMenuOptions() []MenuOption {
	return []MenuOption{
		{Label: "🔑 Configurar API Key", Value: "apikey"},
		{Label: "🤖 Seleccionar Modelo", Value: "model"},
		{Label: "📁 Configurar Carpeta Base", Value: "folder"},
		{Label: "📄 Actualizar Template (CSS)", Value: "css"},
		{Label: "📄 Actualizar Prompt (YAML)", Value: "yaml"},
		{Label: "💰 Actualizar Presupuesto (YAML)", Value: "presupuesto_yaml"},
		{Label: "🖼️  Actualizar Logo", Value: "logo"},
		{Label: "⬅️  Volver", Value: "back"},
	}
}

// ShowMainMenu displays the main menu and returns the selected option
func ShowMainMenu() (string, error) {
	fmt.Println(Banner())

	options := MainMenuOptions()
	opts := make([]huh.Option[string], len(options))
	for i, opt := range options {
		opts[i] = huh.NewOption(opt.Label, opt.Value)
	}

	var selected string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Menú Principal").
				Options(opts...).
				Value(&selected),
		),
	).WithTheme(getTheme())

	if err := form.Run(); err != nil {
		return "", err
	}

	return selected, nil
}

// ShowConfigMenu displays the config menu and returns the selected option
func ShowConfigMenu() (string, error) {
	options := ConfigMenuOptions()
	opts := make([]huh.Option[string], len(options))
	for i, opt := range options {
		opts[i] = huh.NewOption(opt.Label, opt.Value)
	}

	var selected string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Configuración").
				Options(opts...).
				Value(&selected),
		),
	).WithTheme(getTheme())

	if err := form.Run(); err != nil {
		return "", err
	}

	return selected, nil
}

// ShowProjectList displays a list of projects and returns the selected one
func ShowProjectList(projects []string) (string, error) {
	if len(projects) == 0 {
		PrintWarning("No hay proyectos disponibles")
		return "", nil
	}

	opts := make([]huh.Option[string], len(projects)+1)
	for i, proj := range projects {
		opts[i] = huh.NewOption(proj, proj)
	}
	opts[len(projects)] = huh.NewOption("⬅️ Volver", "back")

	var selected string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Seleccionar Proyecto").
				Options(opts...).
				Value(&selected),
		),
	).WithTheme(getTheme())

	if err := form.Run(); err != nil {
		return "", err
	}

	return selected, nil
}

// ShowModelSelector displays a model selector and returns the selected model
func ShowModelSelector(models []string, currentModel string) (string, error) {
	opts := make([]huh.Option[string], len(models))
	for i, model := range models {
		label := model
		if model == currentModel {
			label = model + " (actual)"
		}
		opts[i] = huh.NewOption(label, model)
	}

	var selected string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Seleccionar Modelo de IA").
				Options(opts...).
				Value(&selected),
		),
	).WithTheme(getTheme())

	if err := form.Run(); err != nil {
		return "", err
	}

	return selected, nil
}

// ProposalSummary represents a proposal summary for display
type ProposalSummary struct {
	Project  string
	Title    string
	Date     string
	FilePath string
}

// PresupuestoSummary represents a budget summary for display
type PresupuestoSummary struct {
	Project  string
	Date     string
	FilePath string
}

// ShowProposalSummaries displays a list of proposal summaries
func ShowProposalSummaries(summaries []ProposalSummary) (string, error) {
	if len(summaries) == 0 {
		PrintWarning("No hay propuestas disponibles")
		return "", nil
	}

	opts := make([]huh.Option[string], len(summaries)+1)
	for i, summary := range summaries {
		label := fmt.Sprintf("%s | %s | %s", summary.Project, summary.Title, summary.Date)
		opts[i] = huh.NewOption(label, summary.FilePath)
	}
	opts[len(summaries)] = huh.NewOption("⬅️ Volver", "back")

	var selected string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Propuestas Disponibles").
				Options(opts...).
				Value(&selected),
		),
	).WithTheme(getTheme())

	if err := form.Run(); err != nil {
		return "", err
	}

	return selected, nil
}

// ShowPresupuestoSummaries displays a list of budget summaries
func ShowPresupuestoSummaries(summaries []PresupuestoSummary) (string, error) {
	if len(summaries) == 0 {
		PrintWarning("No hay presupuestos disponibles")
		return "", nil
	}

	opts := make([]huh.Option[string], len(summaries)+1)
	for i, summary := range summaries {
		label := fmt.Sprintf("%s | %s", summary.Project, summary.Date)
		opts[i] = huh.NewOption(label, summary.FilePath)
	}
	opts[len(summaries)] = huh.NewOption("⬅️ Volver", "back")

	var selected string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Presupuestos Disponibles").
				Options(opts...).
				Value(&selected),
		),
	).WithTheme(getTheme())

	if err := form.Run(); err != nil {
		return "", err
	}

	return selected, nil
}

