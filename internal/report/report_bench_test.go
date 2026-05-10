package report

import (
	"fmt"
	"io"
	"testing"
)

func BenchmarkWrite(b *testing.B) {
	counts := make(map[string]int, 10000)
	for i := 0; i < 10000; i++ {
		counts[fmt.Sprintf("name_%d", i)] = i
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := Write(io.Discard, counts, Options{SortByCount: true}); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWrite_TopN(b *testing.B) {
	counts := make(map[string]int, 10000)
	for i := 0; i < 10000; i++ {
		counts[fmt.Sprintf("name_%d", i)] = i
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := Write(io.Discard, counts, Options{SortByCount: true, Top: 10}); err != nil {
			b.Fatal(err)
		}
	}
}
