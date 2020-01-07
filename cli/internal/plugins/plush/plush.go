package plush

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/buildcmd"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/plush"
)

type Buffalo struct{}

var _ plugins.Plugin = Buffalo{}

func (b Buffalo) Name() string {
	return "plush"
}

var _ buildcmd.TemplatesValidator = &Buffalo{}

func (b *Buffalo) ValidateTemplates(root string) error {
	root = filepath.Join(root, "templates")
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
