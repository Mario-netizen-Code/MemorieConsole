//go:build !windows

package app

import (
	"fmt"
	"os"
	"strconv"
)

func openFileShared(path string) (*os.File, error) {
	return os.Open(path)
}

func clearScreen() {
	fmt.Print("\033[2J\033[H")
}

func terminalWidth() int {
	cols := os.Getenv("COLUMNS")
	if cols != "" {
		if n, err := strconv.Atoi(cols); err == nil && n > 20 {
			return n
		}
	}
	return 120
}
