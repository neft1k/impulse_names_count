package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_FromStdin(t *testing.T) {
	var stdout, stderr bytes.Buffer
	in := strings.NewReader("Алёна\nМиша\nАлёна\nДима\n")

	if err := Run(nil, in, &stdout, &stderr); err != nil {
		t.Fatalf("unexpected error: %v (stderr=%s)", err, stderr.String())
	}

	want := "Алёна:2\nДима:1\nМиша:1\n"
	if got := stdout.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRun_FromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "names.txt")
	if err := os.WriteFile(path, []byte("Иван\nИван\nОлег\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if err := Run([]string{path}, nil, &stdout, &stderr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "Иван:2\nОлег:1\n"
	if got := stdout.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRun_AlphaFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer
	in := strings.NewReader("в\nа\nб\nа\n")

	if err := Run([]string{"-alpha"}, in, &stdout, &stderr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "а:2\nб:1\nв:1\n"
	if got := stdout.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRun_TopFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer
	in := strings.NewReader("a\na\na\nb\nb\nc\n")

	if err := Run([]string{"-top", "2"}, in, &stdout, &stderr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "a:3\nb:2\n"
	if got := stdout.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRun_FileNotFound(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := Run([]string{"/nope/does/not/exist.txt"}, nil, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestRun_TooManyArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := Run([]string{"a.txt", "b.txt"}, nil, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error for too many args, got nil")
	}
}
