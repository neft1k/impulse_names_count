package counter

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

type countCase struct {
	name  string
	input string
	want  map[string]int
}

var countCases = []countCase{
	{
		name:  "empty input",
		input: "",
		want:  map[string]int{},
	},
	{
		name:  "task example",
		input: "Алёна\nМиша\nАлёна\nДима\n",
		want:  map[string]int{"Алёна": 2, "Миша": 1, "Дима": 1},
	},
	{
		name:  "single name many times",
		input: "Иван\nИван\nИван\n",
		want:  map[string]int{"Иван": 3},
	},
	{
		name:  "trims surrounding whitespace",
		input: "  Анна  \n\tПётр\t\nАнна\n",
		want:  map[string]int{"Анна": 2, "Пётр": 1},
	},
	{
		name:  "skips empty and whitespace-only lines",
		input: "Лиза\n\n   \nЛиза\n\n",
		want:  map[string]int{"Лиза": 2},
	},
	{
		name:  "case sensitive — different names",
		input: "анна\nАнна\n",
		want:  map[string]int{"анна": 1, "Анна": 1},
	},
	{
		name:  "ё and е are different names",
		input: "Алёна\nАлена\n",
		want:  map[string]int{"Алёна": 1, "Алена": 1},
	},
	{
		name:  "no trailing newline",
		input: "Олег\nОлег",
		want:  map[string]int{"Олег": 2},
	},
}

func TestCount(t *testing.T) {
	for _, tt := range countCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Count(strings.NewReader(tt.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

type errReader struct{ err error }

func (e errReader) Read(_ []byte) (int, error) { return 0, e.err }

func TestCount_ReadError(t *testing.T) {
	sentinel := errors.New("disk on fire")
	_, err := Count(errReader{err: sentinel})
	if !errors.Is(err, sentinel) {
		t.Fatalf("want wrapped sentinel error, got %v", err)
	}
}

func TestCount_LineTooLong(t *testing.T) {
	huge := strings.Repeat("x", maxLineSize+1)
	_, err := Count(strings.NewReader(huge))
	if err == nil {
		t.Fatal("expected error for oversized line, got nil")
	}
}
