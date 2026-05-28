package app

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
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
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	green   = "\033[38;5;82m"
	yellow  = "\033[38;5;220m"
	blue    = "\033[38;5;75m"
	cyan    = "\033[38;5;51m"
	pink    = "\033[38;5;207m"
	gray    = "\033[38;5;245m"
	lGray   = "\033[38;5;250m"
	white   = "\033[38;5;255m"
	bgDark  = "\033[48;5;235m"
	bgRed   = "\033[48;5;52m"
	red     = "\033[38;5;196m"
	cursorH = "\033[?25l"
	cursorS = "\033[?25h"
)

var typeStyle = map[string]struct {
	icon  string
	color string
	bg    string
	label string
}{
	"milestone":      {"★", cyan, "\033[48;5;24m", "MILESTONE"},
	"decision":       {"◈", yellow, "\033[48;5;94m", "DECISION"},
	"modification":   {"◆", green, "\033[48;5;22m", "MODIFY"},
	"implementation": {"●", blue, "\033[48;5;19m", "IMPLEMENT"},
	"verification":   {"▸", pink, "\033[48;5;54m", "VERIFY"},
	"finding":        {"○", white, "\033[48;5;236m", "FINDING"},
	"context":        {"·", gray, "\033[48;5;236m", "CONTEXT"},
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
	bar := style.color + "│" + reset
	indent := " " + bar + " "

	var b strings.Builder

	if e.Important {
		b.WriteString(bold + red + "❗" + reset)
		b.WriteString(" ")
	}

	b.WriteString(" ")
	b.WriteString(style.color + style.icon + reset)
	b.WriteString("  ")
	b.WriteString(style.bg + bold + white + " " + style.label + " " + reset)
	b.WriteString("  ")
	b.WriteString(style.color + formatTime(e.Timestamp) + reset)
	b.WriteString("\n")

	if e.Summary != "" {
		for _, l := range wrap(e.Summary, 0, w) {
			b.WriteString(indent + " " + bold + white + l + reset + "\n")
		}
	}

	if e.Details != "" {
		for _, l := range wrap(e.Details, 0, w) {
			b.WriteString(indent + " " + lGray + l + reset + "\n")
		}
	}

	if len(e.Tags) > 0 {
		b.WriteString(indent + " ")
		for _, tag := range e.Tags {
			b.WriteString(gray + "#" + tag + reset + " ")
		}
		b.WriteString("\n")
	}

	if len(e.Files) > 0 {
		b.WriteString(indent + " ")
		b.WriteString(dim)
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
		b.WriteString(reset + "\n")
	}

	b.WriteString(" " + style.color + "─" + reset + strings.Repeat("─", w-2) + "\n")

	return b.String()
}

func printHeader(files []string) {
	w := terminalWidth()
	title := " MEMORIE CONSOLE "
	left := strings.Repeat("─", (w-utf8.RuneCountInString(title))/2-1)
	right := strings.Repeat("─", w-(w-utf8.RuneCountInString(title))/2-utf8.RuneCountInString(title)-1)

	fmt.Print(reset + bold + cyan + "╭" + left + title + right + "╮" + reset + "\n")

	sub := dim + " " + strings.Join(files, ", ") + reset
	if utf8.RuneCountInString(sub) > w-2 {
		sub = dim + fmt.Sprintf(" %d files", len(files)) + reset
	}
	pad := strings.Repeat(" ", w-utf8.RuneCountInString(sub)-1)
	fmt.Print(cyan + "│" + reset + sub + pad + cyan + "│" + reset + "\n")
	fmt.Print(cyan + "╰" + strings.Repeat("─", w-2) + "╯" + reset + "\n\n")
}

var totalEntries int

func printEntryLine(line string) {
	if strings.HasPrefix(line, "#") {
		fmt.Print(dim + line + reset + "\n")
		return
	}
	if formatted := formatEntry(line); formatted != "" {
		fmt.Print(formatted)
		totalEntries++
	}
}

func cleanupDisplay() {
	fmt.Print(cursorS)
}
