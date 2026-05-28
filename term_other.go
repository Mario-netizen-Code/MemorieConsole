//go:build !windows

package main

import (
	"os"
	"strconv"
)

func terminalWidth() int {
	cols := os.Getenv("COLUMNS")
	if cols != "" {
		if n, err := strconv.Atoi(cols); err == nil && n > 20 {
			return n
		}
	}
	return 120
}
