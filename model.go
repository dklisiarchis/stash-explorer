package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// viewState tracks which level of the UI we're on.
type viewState int

const (
	stashListView viewState = iota
	fileListView
	diffView
)

// applyScope describes what will be applied.
type applyScope int

const (
	applyWholeStash applyScope = iota
	applySingleFile
)

// Async messages for loading data.
type stashesLoadedMsg struct {
	stashes []stashEntry
	err     error
}

type filesLoadedMsg struct {
	files []fileEntry
	err   error
}

type diffLoadedMsg struct {
	diff string
	file string
	err  error
}

type applyResultMsg struct {
	err   error
	label string
}

// model is the top-level Bubble Tea model.
type model struct {
	state  viewState
	width  int
	height int

	// Stash list level
	stashList list.Model
	stashes   []stashEntry

	// File list level
	fileList    list.Model
	files       []fileEntry
	activeStash stashEntry

	// Diff level
	diffViewport viewport.Model
	activeFile   string
	diffContent  string

	// Confirmation
	confirming   bool
	confirmScope applyScope
	confirmRef   string
	confirmFile  string
	confirmLabel string

	// Shared state
	showHelp bool
	err      error
	success  string
	loading  bool
}

const footerHeight = 1

func initialModel() model {
	return model{
		state:   stashListView,
		loading: true,
	}
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		stashes, err := loadStashes()
		return stashesLoadedMsg{stashes: stashes, err: err}
	}
}

func (m model) contentHeight() int {
	// breadcrumb(1) + content + footer(1)
	h := m.height - 2 - footerHeight
	if h < 5 {
		h = 20
	}
	return h
}

func (m model) safeWidth() int {
	if m.width < 20 {
		return 80
	}
	return m.width
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.loading {
			h := m.contentHeight()
			w := m.safeWidth()
			switch m.state {
			case stashListView:
				m.stashList.SetSize(w, h)
			case fileListView:
				m.fileList.SetSize(w, h)
			case diffView:
				m.diffViewport.Width = w
				m.diffViewport.Height = h
			}
		}
		return m, nil

	case stashesLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.stashes = msg.stashes
		m.stashList = newStashList(m.stashes, m.safeWidth(), m.contentHeight())
		return m, nil

	case filesLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.files = msg.files
		m.state = fileListView
		m.fileList = newFileList(m.files, m.safeWidth(), m.contentHeight())
		return m, nil

	case diffLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.state = diffView
		m.activeFile = msg.file
		m.diffContent = msg.diff
		m.diffViewport = newDiffViewport(m.diffContent, m.safeWidth(), m.contentHeight())
		return m, nil

	case applyResultMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.success = msg.label
		}
		return m, nil

	case tea.KeyMsg:
		// Clear success message on any key
		if m.success != "" {
			m.success = ""
		}

		// Handle confirmation dialog first
		if m.confirming {
			return m.updateConfirm(msg)
		}

		// Global keys always work
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if !m.loading {
				if m.state == stashListView && m.stashList.FilterState() == list.Filtering {
					break // let list handle it
				}
				if m.state == fileListView && m.fileList.FilterState() == list.Filtering {
					break
				}
			}
			return m, tea.Quit
		case "?":
			m.showHelp = !m.showHelp
			return m, nil
		case "ctrl+k":
			if m.loading {
				return m, nil
			}
			return m.startConfirm()
		}

		// Don't forward keys while loading
		if m.loading {
			return m, nil
		}

		if m.showHelp {
			m.showHelp = false
			return m, nil
		}

		// Clear errors on Esc
		if msg.String() == "esc" && m.err != nil && m.state != stashListView {
			m.err = nil
		}

		return m.updateForState(msg)
	}

	// Pass other messages to active view (only when not loading)
	if !m.loading {
		return m.updateSubview(msg)
	}
	return m, nil
}

// startConfirm enters the confirmation dialog for the current context.
func (m model) startConfirm() (tea.Model, tea.Cmd) {
	switch m.state {
	case stashListView:
		item, ok := m.stashList.SelectedItem().(stashItem)
		if !ok {
			return m, nil
		}
		m.confirming = true
		m.confirmScope = applyWholeStash
		m.confirmRef = item.entry.ref
		m.confirmLabel = fmt.Sprintf("Apply %s: %s", item.entry.ref, item.entry.message)

	case fileListView:
		m.confirming = true
		m.confirmScope = applyWholeStash
		m.confirmRef = m.activeStash.ref
		m.confirmLabel = fmt.Sprintf("Apply %s: %s", m.activeStash.ref, m.activeStash.message)

	case diffView:
		m.confirming = true
		m.confirmScope = applySingleFile
		m.confirmRef = m.activeStash.ref
		m.confirmFile = m.activeFile
		m.confirmLabel = fmt.Sprintf("Apply %s to %s", m.activeStash.ref, m.activeFile)
	}
	return m, nil
}

func (m model) updateConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		m.confirming = false
		m.loading = true
		m.err = nil
		ref := m.confirmRef
		file := m.confirmFile
		scope := m.confirmScope
		label := m.confirmLabel
		return m, func() tea.Msg {
			var err error
			if scope == applySingleFile {
				err = applyFile(ref, file)
			} else {
				err = applyStash(ref)
			}
			return applyResultMsg{err: err, label: label}
		}
	case "n", "N", "esc":
		m.confirming = false
		return m, nil
	}
	return m, nil
}

func (m model) updateForState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state {
	case stashListView:
		return m.updateStashList(msg)
	case fileListView:
		return m.updateFileList(msg)
	case diffView:
		return m.updateDiffView(msg)
	}
	return m, nil
}

func (m model) updateStashList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if m.stashList.FilterState() == list.Filtering {
			break // let the list handle enter to confirm filter
		}
		item, ok := m.stashList.SelectedItem().(stashItem)
		if !ok {
			return m, nil
		}
		m.activeStash = item.entry
		m.loading = true
		m.err = nil
		ref := item.entry.ref
		return m, func() tea.Msg {
			files, err := loadFiles(ref)
			return filesLoadedMsg{files: files, err: err}
		}
	case "esc":
		if m.stashList.FilterState() == list.Filtering {
			break // let list cancel filter
		}
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.stashList, cmd = m.stashList.Update(msg)
	return m, cmd
}

func (m model) updateFileList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if m.fileList.FilterState() == list.Filtering {
			break
		}
		item, ok := m.fileList.SelectedItem().(fileItem)
		if !ok {
			return m, nil
		}
		m.loading = true
		m.err = nil
		ref := m.activeStash.ref
		file := item.entry.name
		return m, func() tea.Msg {
			diff, err := loadDiff(ref, file)
			return diffLoadedMsg{diff: diff, file: file, err: err}
		}
	case "esc":
		if m.fileList.FilterState() == list.Filtering {
			break
		}
		m.state = stashListView
		m.err = nil
		return m, nil
	}

	var cmd tea.Cmd
	m.fileList, cmd = m.fileList.Update(msg)
	return m, cmd
}

func (m model) updateDiffView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = fileListView
		m.err = nil
		return m, nil
	}

	var cmd tea.Cmd
	m.diffViewport, cmd = m.diffViewport.Update(msg)
	return m, cmd
}

func (m model) updateSubview(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case stashListView:
		m.stashList, cmd = m.stashList.Update(msg)
	case fileListView:
		m.fileList, cmd = m.fileList.Update(msg)
	case diffView:
		m.diffViewport, cmd = m.diffViewport.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	if m.showHelp {
		return renderHelp(m.width, m.height)
	}

	if m.confirming {
		return m.viewConfirm()
	}

	if m.loading {
		return breadcrumbStyle.Render("Loading…")
	}

	if m.err != nil && m.state == stashListView {
		return errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	var content string
	switch m.state {
	case stashListView:
		content = m.viewStashList()
	case fileListView:
		content = m.viewFileList()
	case diffView:
		content = m.viewDiff()
	}

	return content + "\n" + m.viewFooter()
}

func (m model) viewConfirm() string {
	title := confirmTitleStyle.Render("Confirm Apply")
	desc := "\n\n" + m.confirmLabel
	if m.confirmScope == applySingleFile {
		desc += "\n\nThis will restore this file from the stash into your working tree."
	} else {
		desc += "\n\nThis will apply all changes from the stash to your working tree."
	}
	hint := "\n\n" + confirmHintStyle.Render("y to confirm / n or Esc to cancel")

	box := confirmStyle.
		Width(min(60, m.width-4)).
		Render(title + desc + hint)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}

func (m model) viewFooter() string {
	var left string

	if m.success != "" {
		left = successStyle.Render("Applied: " + m.success)
	} else if m.err != nil {
		left = errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	applyLabel := "Apply stash"
	if m.state == diffView {
		applyLabel = "Apply file"
	}

	right := footerKeyStyle.Render("^K") + " " + footerDescStyle.Render(applyLabel) +
		footerKeyStyle.Render("?") + " " + footerDescStyle.Render("Help")

	if m.state == diffView {
		scrollPct := fmt.Sprintf(" %3.f%%", m.diffViewport.ScrollPercent()*100)
		right = statusBarStyle.Render(scrollPct) + "  " + right
	}

	availWidth := m.width - lipgloss.Width(right)
	if availWidth < 0 {
		availWidth = 0
	}

	if left != "" {
		left = lipgloss.NewStyle().Width(availWidth).Render(left)
		return left + right
	}

	return lipgloss.PlaceHorizontal(m.width, lipgloss.Right, right)
}

func (m model) breadcrumb() string {
	switch m.state {
	case stashListView:
		return breadcrumbStyle.Render("Stashes")
	case fileListView:
		stashLabel := fmt.Sprintf("%s: %s", m.activeStash.ref, m.activeStash.message)
		return breadcrumbStyle.Render("Stashes") +
			breadcrumbSep.String() +
			breadcrumbStyle.Render(truncate(stashLabel, 40))
	case diffView:
		stashLabel := fmt.Sprintf("%s: %s", m.activeStash.ref, m.activeStash.message)
		return breadcrumbStyle.Render("Stashes") +
			breadcrumbSep.String() +
			breadcrumbStyle.Render(truncate(stashLabel, 30)) +
			breadcrumbSep.String() +
			breadcrumbStyle.Render(truncate(m.activeFile, 30))
	}
	return ""
}

func (m model) viewStashList() string {
	return m.breadcrumb() + "\n" + m.stashList.View()
}

func (m model) viewFileList() string {
	header := m.breadcrumb()
	return header + "\n" + m.fileList.View()
}

func (m model) viewDiff() string {
	header := m.breadcrumb()
	return header + "\n" + m.diffViewport.View()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
