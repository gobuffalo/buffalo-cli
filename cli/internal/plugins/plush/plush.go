package plush

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/plush"
)

type Buffalo struct {
	stdin         io.Reader
	stdout        io.Writer
	stderr        io.Writer
	TemplatesPath string
}

func (b Buffalo) Name() string {
	return "plush"
}

func (b *Buffalo) Stderr() io.Writer {
	if b.stderr != nil {
		return b.stderr
	}
	return os.Stderr
}

func (b *Buffalo) SetStderr(w io.Writer) {
	b.stderr = w
}

func (b *Buffalo) Stdin() io.Reader {
	if b.stdin != nil {
		return b.stdin
	}
	return os.Stdin
}

func (b *Buffalo) SetStdin(r io.Reader) {
	b.stdin = r
}

func (b *Buffalo) Stdout() io.Writer {
	if b.stdout != nil {
		return b.stdout
	}
	return os.Stdout
}

func (b *Buffalo) SetStdout(w io.Writer) {
	b.stdout = w
}

func (b *Buffalo) ValidateTemplates(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		base := filepath.Base(path)
		if !strings.Contains(base, ".plush") {
			return nil
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		if _, err = plush.Parse(string(b)); err != nil {
			return fmt.Errorf("could not parse %s: %v", path, err)
		}
		return nil
	})
}
