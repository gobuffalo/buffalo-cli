package templates

import (
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/build"
	"github.com/gobuffalo/plugins"
)

var _ plugins.Plugin = Validator{}
var _ build.BeforeBuilder = &Validator{}

type Validator struct{}

func (t *Validator) BeforeBuild(ctx context.Context, root string, args []string) error {
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

func (t Validator) PluginName() string {
	return "golang/templates"
}
