package ui

import "github.com/charmbracelet/lipgloss"

// Color palette - azul, cyan, rojo, naranja
var (
	// Primary colors
	ColorBlue      = lipgloss.Color("#476FD6")
	ColorDarkBlue  = lipgloss.Color("#1C3F99")
	ColorSkyBlue   = lipgloss.Color("#87CEEB")
	ColorCyan      = lipgloss.Color("#39")
	ColorRed       = lipgloss.Color("#204")
	ColorOrange    = lipgloss.Color("#214")
	ColorGreen     = lipgloss.Color("#42")
	ColorWhite     = lipgloss.Color("#FFFFFF")
	ColorGray      = lipgloss.Color("#666666")

	// Styles
	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorBlue).
			Bold(true).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorSkyBlue).
			Italic(true)

	InfoStyle = lipgloss.NewStyle().
			Foreground(ColorCyan).
			Bold(false)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorGreen).
			Bold(false)

	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorOrange).
			Bold(false)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorRed).
			Bold(true)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorBlue).
			Padding(1, 2)

	HeaderStyle = lipgloss.NewStyle().
			Background(ColorBlue).
			Foreground(ColorWhite).
			Bold(true).
			Padding(0, 2)

	MenuItemStyle = lipgloss.NewStyle().
			Foreground(ColorSkyBlue).
			PaddingLeft(2)

	SelectedMenuItemStyle = lipgloss.NewStyle().
				Foreground(ColorBlue).
				Bold(true).
				PaddingLeft(2)

	PromptStyle = lipgloss.NewStyle().
			Foreground(ColorCyan).
			Bold(true)
)

// Banner returns the application banner
func Banner() string {
	banner := `
╔══════════════════════════════════════════════════════╗
║                    ORGMPROP CLI                      ║
║           Generador de Propuestas                    ║
╚══════════════════════════════════════════════════════╝`
	return TitleStyle.Render(banner)
}

