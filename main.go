package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	flag.StringVar(&repoDir, "C", "", "Run as if git was started in this directory")
	flag.Parse()

	if !isGitRepo() {
		fmt.Fprintln(os.Stderr, "Error: not a git repository (or any parent up to mount point /)")
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
