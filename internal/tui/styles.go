package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	PrimaryPink = lipgloss.Color("#FF69B4") // Hot Pink
	SoftPink    = lipgloss.Color("#FFB6C1") // Light Pink
	AccentPink  = lipgloss.Color("#FF1493") // Deep Pink
	White       = lipgloss.Color("#FFFFFF")
	Gray        = lipgloss.Color("#808080")

	// Styles
	TitleStyle = lipgloss.NewStyle().
			Foreground(AccentPink).
			Bold(true).
			Padding(0, 1)

	BorderPink = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryPink).
			Padding(1, 2)

	TimerStyle = lipgloss.NewStyle().
			Foreground(White).
			Background(AccentPink).
			Padding(0, 1).
			Bold(true)

	PlayerNameStyle = lipgloss.NewStyle().
			Foreground(SoftPink).
			Bold(true)

	ChipStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")) // Gold

	WinStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)

	LoseStyle = lipgloss.NewStyle().
			Foreground(Gray)

	BetTypeStyle = lipgloss.NewStyle().
			Foreground(White).
			Background(PrimaryPink).
			Padding(0, 1)

	SloganStyle = lipgloss.NewStyle().
			Foreground(SoftPink).
			Italic(true)
)
