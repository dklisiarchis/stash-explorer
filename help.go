package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type helpBinding struct {
	key  string
	desc string
}

var helpBindings = []helpBinding{
	{"Enter", "Drill into stash / file"},
	{"Esc", "Go back / quit"},
	{"q / Ctrl+C", "Quit"},
	{"/", "Filter list"},
	{"j/k / ↑/↓", "Navigate"},
	{"PgUp/PgDn", "Scroll diff"},
	{"Ctrl+K", "Apply stash / file"},
	{"?", "Toggle this help"},
}

// renderHelp returns the styled help overlay content.
func renderHelp(width, height int) string {
	var rows []string
	for _, b := range helpBindings {
		row := helpKeyStyle.Render(b.key) + helpDescStyle.Render(b.desc)
		rows = append(rows, row)
	}

	content := "Keyboard Shortcuts\n\n" + strings.Join(rows, "\n")

	overlay := helpStyle.
		Width(min(50, width-4)).
		Render(content)

	// Center the overlay
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, overlay)
}
