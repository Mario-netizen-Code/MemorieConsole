package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"
)

type Entry struct {
	ID        string   `json:"id"`
	Timestamp string   `json:"timestamp"`
	SessionID int      `json:"session_id"`
	Scope     string   `json:"scope"`
	Important bool     `json:"important"`
	Type      string   `json:"type"`
	Summary   string   `json:"summary"`
	Files     []string `json:"files"`
	Branch    string   `json:"branch"`
	Tags      []string `json:"tags"`
	Details   string   `json:"details"`
}

const (
	Reset   = "\033[0m"
	Bold    = "\033[1m"
	Dim     = "\033[2m"
	Green   = "\033[38;5;82m"
	Yellow  = "\033[38;5;220m"
	Blue    = "\033[38;5;75m"
	Cyan    = "\033[38;5;51m"
	Pink    = "\033[38;5;207m"
	Gray    = "\033[38;5;245m"
	LGray   = "\033[38;5;250m"
	White   = "\033[38;5;255m"
	BgDark  = "\033[48;5;235m"
	BgRed   = "\033[48;5;52m"
	Red     = "\033[38;5;196m"
	CursorH = "\033[?25l"
	CursorS = "\033[?25h"
)

var typeStyle = map[string]struct {
	icon  string
	color string
	bg    string
	label string
}{
	"milestone":      {"★", Cyan, "\033[48;5;24m", "MILESTONE"},
	"decision":       {"◈", Yellow, "\033[48;5;94m", "DECISION"},
	"modification":   {"◆", Green, "\033[48;5;22m", "MODIFY"},
	"implementation": {"●", Blue, "\033[48;5;19m", "IMPLEMENT"},
	"verification":   {"▸", Pink, "\033[48;5;54m", "VERIFY"},
	"finding":        {"○", White, "\033[48;5;236m", "FINDING"},
	"context":        {"·", Gray, "\033[48;5;236m", "CONTEXT"},
}

type fileTail struct {
	file    *os.File
	lastPos int64
	path    string
}

func formatTime(ts string) string {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		if len(ts) >= 19 {
			return ts[11:19]
		}
		return ts
	}
	return t.Local().Format("15:04:05")
}

func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n-1]) + "…"
}

func wrap(s string, indent int, width int) []string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return nil
	}
	var lines []string
	line := strings.Repeat(" ", indent)
	for _, w := range words {
		if utf8.RuneCountInString(line)+utf8.RuneCountInString(w)+1 > width {
			lines = append(lines, line)
			line = strings.Repeat(" ", indent) + w
		} else if line == strings.Repeat(" ", indent) {
			line += w
		} else {
			line += " " + w
		}
	}
	if line != strings.Repeat(" ", indent) {
		lines = append(lines, line)
	}
	return lines
}

func formatEntry(line string) string {
	var e Entry
	if err := json.Unmarshal([]byte(line), &e); err != nil {
		return ""
	}

	style, ok := typeStyle[e.Type]
	if !ok {
		style = typeStyle["finding"]
	}

	w := terminalWidth() - 5
	bar := style.color + "│" + Reset
	indent := " " + bar + " "

	var b strings.Builder

	if e.Important {
		b.WriteString(Bold + Red + "❗" + Reset)
		b.WriteString(" ")
	}

	b.WriteString(" ")
	b.WriteString(style.color + style.icon + Reset)
	b.WriteString("  ")
	b.WriteString(style.bg + Bold + White + " " + style.label + " " + Reset)
	b.WriteString("  ")
	b.WriteString(style.color + formatTime(e.Timestamp) + Reset)
	b.WriteString("\n")

	if e.Summary != "" {
		for _, l := range wrap(e.Summary, 0, w) {
			b.WriteString(indent + " " + Bold + White + l + Reset + "\n")
		}
	}

	if e.Details != "" {
		for _, l := range wrap(e.Details, 0, w) {
			b.WriteString(indent + " " + LGray + l + Reset + "\n")
		}
	}

	if len(e.Tags) > 0 {
		b.WriteString(indent + " ")
		for _, tag := range e.Tags {
			b.WriteString(Gray + "#" + tag + Reset + " ")
		}
		b.WriteString("\n")
	}

	if len(e.Files) > 0 {
		b.WriteString(indent + " ")
		b.WriteString(Dim)
		for i, fp := range e.Files {
			if i > 0 {
				b.WriteString(", ")
			}
			icon := "📄"
			if info, err := os.Stat(fp); err == nil && info.IsDir() {
				icon = "📁"
			}
			b.WriteString(icon + " " + fp)
		}
		b.WriteString(Reset + "\n")
	}

	b.WriteString(" " + style.color + "─" + Reset + strings.Repeat("─", w-2) + "\n")

	return b.String()
}

var totalEntries int

func printHeader(files []string) {
	w := terminalWidth()
	title := " MEMORIE CONSOLE "
	left := strings.Repeat("─", (w-utf8.RuneCountInString(title))/2-1)
	right := strings.Repeat("─", w-(w-utf8.RuneCountInString(title))/2-utf8.RuneCountInString(title)-1)

	fmt.Print(Reset + Bold + Cyan + "╭" + left + title + right + "╮" + Reset + "\n")

	sub := Dim + " " + strings.Join(files, ", ") + Reset
	if utf8.RuneCountInString(sub) > w-2 {
		sub = Dim + fmt.Sprintf(" %d files", len(files)) + Reset
	}
	pad := strings.Repeat(" ", w-utf8.RuneCountInString(sub)-1)
	fmt.Print(Cyan + "│" + Reset + sub + pad + Cyan + "│" + Reset + "\n")
	fmt.Print(Cyan + "╰" + strings.Repeat("─", w-2) + "╯" + Reset + "\n\n")
}

func printEntryLine(line string) {
	if strings.HasPrefix(line, "#") {
		fmt.Print(Dim + line + Reset + "\n")
		return
	}
	if formatted := formatEntry(line); formatted != "" {
		fmt.Print(formatted)
		totalEntries++
	}
}

func readAllLines(f *os.File, lastPos int64) ([]string, int64) {
	if lastPos > 0 {
		_, err := f.Seek(lastPos, 0)
		if err != nil {
			return nil, lastPos
		}
	}

	var lines []string
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 64*1024), 1024*1024)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	pos, _ := f.Seek(0, 1)
	return lines, pos
}

func getMemFiles(dir string) []string {
	matches, err := filepath.Glob(filepath.Join(dir, "*.txt"))
	if err != nil || len(matches) == 0 {
		return nil
	}
	sort.Strings(matches)
	return matches
}

func cleanup() {
	fmt.Print(CursorS)
}

func main() {
	var paths []string

	if len(os.Args) > 1 {
		paths = os.Args[1:]
	} else {
		memFiles := getMemFiles(".memories")
		if len(memFiles) > 0 {
			paths = memFiles
		} else {
			paths = []string{"convertfiletomarkdown_2026-05-27.txt"}
		}
	}

	var tails []*fileTail
	var sourceNames []string

	for _, p := range paths {
		f, err := openFileShared(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening %s: %v\n", p, err)
			continue
		}
		name := filepath.Base(p)
		tails = append(tails, &fileTail{file: f, path: name})
		sourceNames = append(sourceNames, name)
	}

	if len(tails) == 0 {
		fmt.Fprintln(os.Stderr, "No files to read")
		os.Exit(1)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	clearScreen()
	fmt.Print(CursorH)
	defer cleanup()

	printHeader(sourceNames)

	for _, t := range tails {
		lines, pos := readAllLines(t.file, 0)
		t.lastPos = pos
		for _, line := range lines {
			printEntryLine(line)
		}
	}

	fmt.Print(Dim + " ── Loaded " + Bold + fmt.Sprint(totalEntries) + Reset + Dim + " entries from " + Bold + fmt.Sprint(len(tails)) + Reset + Dim + " file")
	if len(tails) != 1 {
		fmt.Print("s")
	}
	fmt.Print(" ──" + Reset + "\n\n")

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	knownPaths := make(map[string]bool)
	for _, t := range tails {
		knownPaths[t.path] = true
	}

	for {
		select {
		case <-sigCh:
			for _, t := range tails {
				t.file.Close()
			}
			return
		case <-ticker.C:
		for _, t := range tails {
			stat, err := t.file.Stat()
			if err != nil {
				continue
			}
			if stat.Size() <= t.lastPos {
				continue
			}

			lines, pos := readAllLines(t.file, t.lastPos)
			t.lastPos = pos
			for _, line := range lines {
				printEntryLine(line)
			}
		}

		if len(os.Args) <= 1 {
			newFiles := getMemFiles(".memories")
			for _, fp := range newFiles {
				name := filepath.Base(fp)
				if knownPaths[name] {
					continue
				}
				f, err := openFileShared(fp)
				if err != nil {
					continue
				}
				knownPaths[name] = true
				tails = append(tails, &fileTail{file: f, path: name})
				lines, pos := readAllLines(f, 0)
				fi := tails[len(tails)-1]
				fi.lastPos = pos
				fmt.Print(Dim + "── New file: " + name + " ──" + Reset + "\n")
				for _, line := range lines {
					printEntryLine(line)
				}
			}
		}
	}
}
}
