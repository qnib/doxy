package proxy

import (
	"bufio"
	"strings"
	"io"
)

func ReadPatterns(reader io.Reader) (patterns []string, err error) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if ! strings.HasPrefix(line, "#") {
			patterns = append(patterns, line)
		}
	}
	return
}
