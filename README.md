# stash-explorer

A lightweight TUI tool for interactively exploring git stashes. Browse stashes, drill into changed files, and view colorized diffs — all from the terminal.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)

## Features

- **Three-level navigation**: Stash list → File list → Diff view
- **Colorized diffs**: Green additions, red deletions, cyan hunk headers
- **Line stats**: See `+N -M` counts per file at a glance
- **Fuzzy filtering**: Press `/` to search stashes or files
- **Apply stashes**: Apply a whole stash or a single file with `Ctrl+K`
- **Confirmation prompts**: Always confirms before modifying your working tree
- **Mouse scroll**: Scroll through diffs with your mouse wheel
- **Breadcrumb navigation**: Always know where you are

## Install

```bash
go install github.com/dklisiarchis/stash-explorer@latest
```

Or build from source:

```bash
git clone https://github.com/dklisiarchis/stash-explorer.git
cd stash-explorer
go build
```

## Usage

```bash
# Run in the current directory
stash-explorer

# Run against a different repo
stash-explorer -C /path/to/repo
```

## Key Bindings

| Key | Action |
|---|---|
| `Enter` | Drill into stash / file |
| `Esc` | Go back one level (quit from top) |
| `q` / `Ctrl+C` | Quit |
| `/` | Filter list |
| `j/k` / `↑/↓` | Navigate |
| `PgUp` / `PgDn` | Scroll diff |
| `Ctrl+K` | Apply stash or file |
| `?` | Toggle help |

## Built With

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) — TUI components (list, viewport)
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — Style definitions

## License

MIT
