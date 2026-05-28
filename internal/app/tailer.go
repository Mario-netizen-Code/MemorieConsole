package app

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
)

type fileTail struct {
	file    *os.File
	lastPos int64
	path    string
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
