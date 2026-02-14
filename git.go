package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// repoDir is the directory to run git commands in. Set via -C flag.
var repoDir string

// runGit executes a git command and returns its trimmed stdout.
func runGit(args ...string) (string, error) {
	if repoDir != "" {
		args = append([]string{"-C", repoDir}, args...)
	}
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("git %s: %s", args[0], strings.TrimSpace(string(exitErr.Stderr)))
		}
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// isGitRepo checks whether the current (or specified) directory is inside a git repo.
func isGitRepo() bool {
	_, err := runGit("rev-parse", "--git-dir")
	return err == nil
}

// stashEntry represents a single git stash.
type stashEntry struct {
	index   int
	ref     string // e.g. stash@{0}
	branch  string
	message string
}

// parseStashList parses the output of `git stash list`.
// Each line looks like: stash@{0}: On main: fix login bug
// or: stash@{0}: WIP on main: abc1234 commit message
func parseStashList(raw string) []stashEntry {
	if raw == "" {
		return nil
	}
	lines := strings.Split(raw, "\n")
	entries := make([]stashEntry, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		e := stashEntry{index: i, ref: fmt.Sprintf("stash@{%d}", i)}

		// Split on first ": " to get ref and rest
		parts := strings.SplitN(line, ": ", 3)
		if len(parts) >= 3 {
			// parts[1] is like "On main" or "WIP on main"
			branchPart := parts[1]
			branchPart = strings.TrimPrefix(branchPart, "WIP on ")
			branchPart = strings.TrimPrefix(branchPart, "On ")
			e.branch = branchPart
			e.message = parts[2]
		} else if len(parts) == 2 {
			e.message = parts[1]
		} else {
			e.message = line
		}

		entries = append(entries, e)
	}
	return entries
}

// fileEntry represents a file changed in a stash.
type fileEntry struct {
	status  string // A, M, D, R, etc.
	name    string
	added   int // lines added
	removed int // lines removed
}

// parseFileList parses the output of `git stash show --name-status`.
// Each line is tab-delimited: M\tfile.go or R100\told.go\tnew.go
func parseFileList(raw string) []fileEntry {
	if raw == "" {
		return nil
	}
	lines := strings.Split(raw, "\n")
	entries := make([]fileEntry, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) < 2 {
			continue
		}
		status := parts[0]
		name := parts[1]

		// Normalize rename status (R100 -> R)
		if strings.HasPrefix(status, "R") {
			status = "R"
			if len(parts) >= 3 {
				name = parts[1] + " -> " + parts[2]
			}
		}

		entries = append(entries, fileEntry{status: status, name: name})
	}
	return entries
}

// loadStashes fetches and parses all stashes.
func loadStashes() ([]stashEntry, error) {
	out, err := runGit("stash", "list")
	if err != nil {
		return nil, err
	}
	return parseStashList(out), nil
}

// parseNumstat parses `git stash show --numstat` output.
// Each line: "10\t5\tfile.go" (added, removed, filename).
// Binary files show "-\t-\tfile".
func parseNumstat(raw string) map[string][2]int {
	result := make(map[string][2]int)
	if raw == "" {
		return result
	}
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) < 3 {
			continue
		}
		added, _ := strconv.Atoi(parts[0])   // "-" for binary → 0
		removed, _ := strconv.Atoi(parts[1]) // "-" for binary → 0
		// For renames, parts[2] may contain "{old => new}" syntax
		name := parts[2]
		result[name] = [2]int{added, removed}
	}
	return result
}

// loadFiles fetches the list of changed files for a stash.
func loadFiles(ref string) ([]fileEntry, error) {
	out, err := runGit("stash", "show", "--name-status", ref)
	if err != nil {
		return nil, err
	}
	entries := parseFileList(out)

	// Get line stats
	numOut, err := runGit("stash", "show", "--numstat", ref)
	if err == nil {
		stats := parseNumstat(numOut)
		for i := range entries {
			name := entries[i].name
			// For renames, try the new name
			if idx := strings.Index(name, " -> "); idx != -1 {
				name = name[idx+4:]
			}
			if s, ok := stats[name]; ok {
				entries[i].added = s[0]
				entries[i].removed = s[1]
			}
		}
	}

	return entries, nil
}

// loadDiff fetches the diff for a specific file in a stash.
func loadDiff(ref, file string) (string, error) {
	// For renamed files, extract the new filename
	if idx := strings.Index(file, " -> "); idx != -1 {
		file = file[idx+4:]
	}
	return runGit("diff", ref+"^", ref, "--", file)
}

// applyStash applies an entire stash to the working tree.
func applyStash(ref string) error {
	_, err := runGit("stash", "apply", ref)
	return err
}

// applyFile restores a single file from a stash into the working tree.
func applyFile(ref, file string) error {
	if idx := strings.Index(file, " -> "); idx != -1 {
		file = file[idx+4:]
	}
	_, err := runGit("checkout", ref, "--", file)
	return err
}
