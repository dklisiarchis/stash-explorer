package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// fileItem wraps fileEntry to implement bubbles list.Item.
type fileItem struct {
	entry fileEntry
}

func (i fileItem) FilterValue() string {
	return i.entry.name
}

// statusIcon returns the styled status indicator for a file.
func statusIcon(status string) string {
	switch status {
	case "A":
		return statusAdded.String()
	case "M":
		return statusModified.String()
	case "D":
		return statusDeleted.String()
	case "R":
		return statusRenamed.String()
	default:
		return statusModified.String()
	}
}

// fileDelegate renders a file item in the list.
type fileDelegate struct{}

func (d fileDelegate) Height() int                             { return 1 }
func (d fileDelegate) Spacing() int                            { return 0 }
func (d fileDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d fileDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	fi, ok := item.(fileItem)
	if !ok {
		return
	}

	icon := statusIcon(fi.entry.status)
	name := fi.entry.name

	// Format line stats: +10 -5
	stats := ""
	if fi.entry.added > 0 || fi.entry.removed > 0 {
		stats = " " +
			diffAddStyle.Render(fmt.Sprintf("+%d", fi.entry.added)) +
			" " +
			diffDelStyle.Render(fmt.Sprintf("-%d", fi.entry.removed))
	}

	maxWidth := m.Width() - 6 - 16 // leave room for stats
	if maxWidth < 20 {
		maxWidth = 20
	}
	if len(name) > maxWidth {
		name = name[:maxWidth-1] + "…"
	}

	cursor := "  "
	if index == m.Index() {
		cursor = "> "
		name = breadcrumbStyle.Render(name)
	}

	fmt.Fprintf(w, "%s%s %s%s", cursor, icon, name, stats)
}

// newFileList creates a configured list for file entries.
func newFileList(entries []fileEntry, width, height int) list.Model {
	items := make([]list.Item, len(entries))
	for i, e := range entries {
		items[i] = fileItem{entry: e}
	}

	l := list.New(items, fileDelegate{}, width, height)
	l.Title = "Changed Files"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)

	return l
}
