package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// stashItem wraps stashEntry to implement bubbles list.Item.
type stashItem struct {
	entry stashEntry
}

func (i stashItem) FilterValue() string {
	return i.entry.message + " " + i.entry.branch
}

// stashDelegate renders a stash item in the list.
type stashDelegate struct{}

func (d stashDelegate) Height() int                             { return 2 }
func (d stashDelegate) Spacing() int                            { return 0 }
func (d stashDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d stashDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	si, ok := item.(stashItem)
	if !ok {
		return
	}

	ref := si.entry.ref
	msg := si.entry.message
	branch := si.entry.branch

	// Truncate message if needed
	maxWidth := m.Width() - 4
	if maxWidth < 20 {
		maxWidth = 20
	}

	title := fmt.Sprintf("%s: %s", ref, msg)
	if len(title) > maxWidth {
		title = title[:maxWidth-1] + "…"
	}

	subtitle := fmt.Sprintf("  on %s", branch)
	if len(subtitle) > maxWidth {
		subtitle = subtitle[:maxWidth-1] + "…"
	}

	cursor := "  "
	if index == m.Index() {
		cursor = "> "
		title = lipglossSelectedTitle(title)
		subtitle = lipglossSelectedSubtitle(subtitle)
	} else {
		title = lipglossNormalTitle(title)
		subtitle = lipglossNormalSubtitle(subtitle)
	}

	fmt.Fprint(w, cursor+title+"\n"+strings.Repeat(" ", 2)+subtitle)
}

func lipglossSelectedTitle(s string) string {
	return breadcrumbStyle.Render(s)
}

func lipglossSelectedSubtitle(s string) string {
	return statusBarStyle.Render(s)
}

func lipglossNormalTitle(s string) string {
	return s
}

func lipglossNormalSubtitle(s string) string {
	return statusBarStyle.Render(s)
}

// newStashList creates a configured list for stash entries.
func newStashList(entries []stashEntry, width, height int) list.Model {
	items := make([]list.Item, len(entries))
	for i, e := range entries {
		items[i] = stashItem{entry: e}
	}

	l := list.New(items, stashDelegate{}, width, height)
	l.Title = "Stashes"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)

	return l
}
