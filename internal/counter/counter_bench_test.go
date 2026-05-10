package counter

import (
	"strings"
	"testing"
)

func BenchmarkCount(b *testing.B) {
	names := []string{"Алёна", "Миша", "Дима", "Иван", "Пётр", "Анна", "Олег", "Лиза"}
	var sb strings.Builder
	for i := 0; i < 100000; i++ {
		sb.WriteString(names[i%len(names)])
		sb.WriteByte('\n')
	}
	data := sb.String()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := Count(strings.NewReader(data)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCount_HighCardinality(b *testing.B) {
	var sb strings.Builder
	for i := 0; i < 100000; i++ {
		sb.WriteString("name_")
		sb.WriteString(strings.Repeat("x", i%50+1))
		sb.WriteByte('\n')
	}
	data := sb.String()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := Count(strings.NewReader(data)); err != nil {
			b.Fatal(err)
		}
	}
}
