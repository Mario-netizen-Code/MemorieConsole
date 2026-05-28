package app

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func Run() {
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
	fmt.Print(cursorH)
	defer cleanupDisplay()

	printHeader(sourceNames)

	for _, t := range tails {
		lines, pos := readAllLines(t.file, 0)
		t.lastPos = pos
		for _, line := range lines {
			printEntryLine(line)
		}
	}

	fmt.Print(dim + " ── Loaded " + bold + fmt.Sprint(totalEntries) + reset + dim + " entries from " + bold + fmt.Sprint(len(tails)) + reset + dim + " file")
	if len(tails) != 1 {
		fmt.Print("s")
	}
	fmt.Print(" ──" + reset + "\n\n")

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
					fmt.Print(dim + "── New file: " + name + " ──" + reset + "\n")
					for _, line := range lines {
						printEntryLine(line)
					}
				}
			}
		}
	}
}
