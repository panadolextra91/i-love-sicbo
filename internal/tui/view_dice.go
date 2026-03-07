package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var diceFrames = [][]string{
	{
		"   _______  ",
		"  /       \\ ",
		" /  o   o  \\",
		" |    o    |",
		" \\  o   o  /",
		"  \\_______/ ",
	},
	{
		"   _______  ",
		"  /       \\ ",
		" /  o       \\",
		" |    o    |",
		" \\       o  /",
		"  \\_______/ ",
	},
	{
		"   _______  ",
		"  /       \\ ",
		" /         \\",
		" |  o   o  |",
		" \\         /",
		"  \\_______/ ",
	},
}

func (m Model) renderDice() string {
	dice1 := m.Dice[0]
	dice2 := m.Dice[1]
	dice3 := m.Dice[2]

	if m.State == StateRolling {
		// Mock rolling animation
		// In real impl, we should use a random index or time-based
		return lipgloss.JoinHorizontal(lipgloss.Top,
			renderSingleDice(1),
			renderSingleDice(2),
			renderSingleDice(3),
		)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top,
		renderSingleDice(dice1),
		renderSingleDice(dice2),
		renderSingleDice(dice3),
	)
}

func renderSingleDice(n int) string {
	// Simple ASCII dice faces 1-6
	var lines []string
	switch n {
	case 1:
		lines = []string{
			"   _______   ",
			"  /       \\  ",
			" /         \\ ",
			" |    o    | ",
			" \\         / ",
			"  \\_______/  ",
		}
	case 2:
		lines = []string{
			"   _______   ",
			"  /       \\  ",
			" /  o       \\ ",
			" |         | ",
			" \\       o  / ",
			"  \\_______/  ",
		}
	case 3:
		lines = []string{
			"   _______   ",
			"  /       \\  ",
			" /  o       \\ ",
			" |    o    | ",
			" \\       o  / ",
			"  \\_______/  ",
		}
	case 4:
		lines = []string{
			"   _______   ",
			"  /       \\  ",
			" /  o   o  \\ ",
			" |         | ",
			" \\  o   o  / ",
			"  \\_______/  ",
		}
	case 5:
		lines = []string{
			"   _______   ",
			"  /       \\  ",
			" /  o   o  \\ ",
			" |    o    | ",
			" \\  o   o  / ",
			"  \\_______/  ",
		}
	case 6:
		lines = []string{
			"   _______   ",
			"  /       \\  ",
			" /  o   o  \\ ",
			" |  o   o  | ",
			" \\  o   o  / ",
			"  \\_______/  ",
		}
	default:
		lines = []string{"?", "?", "?", "?", "?", "?"}
	}

	return lipgloss.NewStyle().Foreground(SoftPink).Render(strings.Join(lines, "\n"))
}
