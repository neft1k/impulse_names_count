package report

import (
	"bufio"
	"fmt"
	"io"
	"sort"
)

type Options struct {
	SortByCount bool
	Top         int
}

func Write(w io.Writer, counts map[string]int, opts Options) error {
	pairs := make([]pair, 0, len(counts))
	for name, n := range counts {
		pairs = append(pairs, pair{name: name, count: n})
	}

	if opts.SortByCount {
		sort.Slice(pairs, func(i, j int) bool {
			if pairs[i].count != pairs[j].count {
				return pairs[i].count > pairs[j].count
			}
			return pairs[i].name < pairs[j].name
		})
	} else {
		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].name < pairs[j].name
		})
	}

	if opts.Top > 0 && opts.Top < len(pairs) {
		pairs = pairs[:opts.Top]
	}

	bw := bufio.NewWriter(w)
	for _, p := range pairs {
		if _, err := fmt.Fprintf(bw, "%s:%d\n", p.name, p.count); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
	}
	return bw.Flush()
}

type pair struct {
	name  string
	count int
}
