package report

import (
	"bytes"
	"strings"
	"testing"
)

type writeCase struct {
	name string
	opts Options
	want string
}

var sampleCounts = map[string]int{
	"Алёна": 2,
	"Миша":  1,
	"Дима":  1,
}

var writeCases = []writeCase{
	{
		name: "sorted by count desc, ties alphabetical",
		opts: Options{SortByCount: true},
		want: "Алёна:2\nДима:1\nМиша:1\n",
	},
	{
		name: "sorted alphabetically",
		opts: Options{SortByCount: false},
		want: "Алёна:2\nДима:1\nМиша:1\n",
	},
	{
		name: "top N truncates",
		opts: Options{SortByCount: true, Top: 2},
		want: "Алёна:2\nДима:1\n",
	},
	{
		name: "top larger than data is a no-op",
		opts: Options{SortByCount: true, Top: 100},
		want: "Алёна:2\nДима:1\nМиша:1\n",
	},
}

func TestWrite(t *testing.T) {
	for _, tt := range writeCases {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := Write(&buf, sampleCounts, tt.opts); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := buf.String(); got != tt.want {
				t.Errorf("got:\n%q\nwant:\n%q", got, tt.want)
			}
		})
	}
}

func TestWrite_EmptyMap(t *testing.T) {
	var buf bytes.Buffer
	if err := Write(&buf, map[string]int{}, Options{SortByCount: true}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got %q", buf.String())
	}
}

func TestWrite_TieOrderingIsStable(t *testing.T) {
	counts := map[string]int{"б": 1, "а": 1, "в": 1}
	var buf bytes.Buffer
	if err := Write(&buf, counts, Options{SortByCount: true}); err != nil {
		t.Fatal(err)
	}
	want := "а:1\nб:1\nв:1\n"
	if got := buf.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

type errWriter struct{}

func (errWriter) Write(_ []byte) (int, error) { return 0, errSentinel }

var errSentinel = &writeErr{"boom"}

type writeErr struct{ msg string }

func (e *writeErr) Error() string { return e.msg }

func TestWrite_PropagatesWriteError(t *testing.T) {
	err := Write(errWriter{}, map[string]int{"a": 1}, Options{})
	if err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("expected wrapped write error, got %v", err)
	}
}
