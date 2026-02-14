package main

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	// Breadcrumb / header
	breadcrumbStyle = lipgloss.NewStyle().
			Foreground(highlight).
			Bold(true).
			PaddingLeft(1)

	breadcrumbSep = lipgloss.NewStyle().
			Foreground(subtle).
			SetString(" > ")

	// Diff colors
	diffAddStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#73F59F"))
	diffDelStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#F5735C"))
	diffHunkStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7EC8E3")).Bold(true)
	diffCtxStyle  = lipgloss.NewStyle()

	// File status indicators
	statusAdded    = lipgloss.NewStyle().Foreground(lipgloss.Color("#73F59F")).SetString("+")
	statusModified = lipgloss.NewStyle().Foreground(lipgloss.Color("#E3D97E")).SetString("~")
	statusDeleted  = lipgloss.NewStyle().Foreground(lipgloss.Color("#F5735C")).SetString("-")
	statusRenamed  = lipgloss.NewStyle().Foreground(lipgloss.Color("#7EC8E3")).SetString("R")

	// Help overlay
	helpStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(highlight).
			Padding(1, 2)

	helpKeyStyle  = lipgloss.NewStyle().Foreground(special).Bold(true).Width(14)
	helpDescStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#555555", Dark: "#AAAAAA"})

	// Error
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F5735C")).
			Bold(true).
			PaddingLeft(1)

	// Status bar
	statusBarStyle = lipgloss.NewStyle().
			Foreground(subtle).
			PaddingLeft(1)

	// Footer bar
	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#555555", Dark: "#AAAAAA"})

	footerKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#F7F7F7", Dark: "#1A1A1A"}).
			Background(lipgloss.AdaptiveColor{Light: "#555555", Dark: "#AAAAAA"}).
			Bold(true).
			Padding(0, 1)

	footerDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#555555", Dark: "#AAAAAA"}).
			PaddingRight(2)

	// Confirmation
	confirmStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#E3D97E")).
			Padding(1, 2)

	confirmTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#E3D97E")).
				Bold(true)

	confirmHintStyle = lipgloss.NewStyle().
				Foreground(subtle)

	// Success
	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#73F59F")).
			Bold(true).
			PaddingLeft(1)
)
