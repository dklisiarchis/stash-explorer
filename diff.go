package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
)

// newDiffViewport creates a configured viewport for displaying a diff.
func newDiffViewport(diff string, width, height int) viewport.Model {
	vp := viewport.New(width, height)
	vp.SetContent(colorizeDiff(diff))
	return vp
}

// colorizeDiff applies lipgloss styles to a unified diff string.
func colorizeDiff(raw string) string {
	lines := strings.Split(raw, "\n")
	styled := make([]string, 0, len(lines))

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "+++ ") || strings.HasPrefix(line, "--- "):
			styled = append(styled, diffHunkStyle.Render(line))
		case strings.HasPrefix(line, "@@"):
			styled = append(styled, diffHunkStyle.Render(line))
		case strings.HasPrefix(line, "+"):
			styled = append(styled, diffAddStyle.Render(line))
		case strings.HasPrefix(line, "-"):
			styled = append(styled, diffDelStyle.Render(line))
		default:
			styled = append(styled, diffCtxStyle.Render(line))
		}
	}

	return strings.Join(styled, "\n")
}
