package git

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

type Buffalo struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func (b *Buffalo) BuildVersion(ctx context.Context, root string) (string, error) {
	now := time.Now()
	version := now.Format(time.RFC3339)
	fmt.Println(">>>TODO cli/internal/plugins/git/git.go:20: version ", version)
	return "", nil
}

// Name is the name of the plugin.
// This will also be used for the cli sub-command
// 	"pop" | "heroku" | "auth" | etc...
func (b *Buffalo) Name() string {
	return "git"
}

func (b *Buffalo) Stdin() io.Reader {
	if b.stdin == nil {
		return os.Stdin
	}
	return b.stdin
}

func (b *Buffalo) Stdout() io.Writer {
	if b.stdout == nil {
		return os.Stdout
	}
	return b.stdout
}

func (b *Buffalo) Stderr() io.Writer {
	if b.stderr == nil {
		return os.Stderr
	}
	return b.stderr
}

func (b *Buffalo) SetStdin(r io.Reader) {
	b.stdin = r
}

func (b *Buffalo) SetStdout(w io.Writer) {
	b.stdout = w
}

func (b *Buffalo) SetStderr(w io.Writer) {
	b.stderr = w
}
