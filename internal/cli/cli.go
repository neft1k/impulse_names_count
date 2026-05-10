package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"impulse_names_count/internal/counter"
	"impulse_names_count/internal/report"
)

func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	opts, err := parseFlags(args, stderr)
	if err != nil {
		return err
	}

	src, closeSrc, err := openInput(opts.files, stdin)
	if err != nil {
		return err
	}
	defer closeSrc()

	counts, err := counter.Count(src)
	if err != nil {
		return err
	}

	return report.Write(stdout, counts, report.Options{
		SortByCount: !opts.alpha,
		Top:         opts.top,
	})
}

type cliOptions struct {
	alpha bool
	top   int
	files []string
}

func parseFlags(args []string, stderr io.Writer) (cliOptions, error) {
	fs := flag.NewFlagSet("count-names", flag.ContinueOnError)
	fs.SetOutput(stderr)

	alpha := fs.Bool("alpha", false, "sort alphabetically")
	top := fs.Int("top", 0, "print only the N most frequent names")

	fs.Usage = func() {
		_, _ = fmt.Fprintf(stderr, "Usage: count-names [flags] <file>\n"+
			"       cat file | count-names [flags]\n\nFlags:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return cliOptions{}, err
	}
	return cliOptions{alpha: *alpha, top: *top, files: fs.Args()}, nil
}

func openInput(files []string, stdin io.Reader) (io.Reader, func(), error) {
	switch len(files) {
	case 0:
		return stdin, func() {}, nil
	case 1:
		f, err := os.Open(files[0])
		if err != nil {
			return nil, nil, fmt.Errorf("open %s: %w", files[0], err)
		}
		return f, func() { _ = f.Close() }, nil
	default:
		return nil, nil, fmt.Errorf("expected at most one file argument, got %d", len(files))
	}
}
