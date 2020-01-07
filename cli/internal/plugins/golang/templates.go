package golang

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/buildcmd"
	"github.com/gobuffalo/buffalo-cli/plugins"
)

type Templates struct{}

var _ buildcmd.TemplatesValidator = &Templates{}

func (t *Templates) ValidateTemplates(root string) error {
	root = filepath.Join(root, "templates")
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		base := filepath.Base(path)
		if !strings.Contains(base, ".tmpl") {
			return nil
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		t := template.New(path)
		if _, err = t.Parse(string(b)); err != nil {
			return fmt.Errorf("could not parse %s: %v", path, err)
		}
		return nil
	})
}

var _ plugins.Plugin = Templates{}

func (t Templates) Name() string {
	return "templates"
}
