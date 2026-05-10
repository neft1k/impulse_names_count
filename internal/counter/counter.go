package counter

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const maxLineSize = 1024 * 1024

func Count(r io.Reader) (map[string]int, error) {
	counts := make(map[string]int)

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 64*1024), maxLineSize)

	for scanner.Scan() {
		name := strings.TrimSpace(scanner.Text())
		if name == "" {
			continue
		}
		counts[name]++
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read input: %w", err)
	}
	return counts, nil
}
