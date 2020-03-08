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

<<<<<<< HEAD:cli/internal/plugins/golang/templates.go
var _ plugins.Plugin = Templater{}
var _ build.BeforeBuilder = &Templater{}
=======
var _ plugins.Plugin = Validator{}
var _ build.TemplatesValidator = &Validator{}
>>>>>>> tweedy:cli/internal/plugins/golang/templates/validator.go

type Validator struct{}

<<<<<<< HEAD:cli/internal/plugins/golang/templates.go
func (t *Templater) BeforeBuild(ctx context.Context, root string, args []string) error {
=======
func (t *Validator) ValidateTemplates(root string) error {
>>>>>>> tweedy:cli/internal/plugins/golang/templates/validator.go
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

<<<<<<< HEAD:cli/internal/plugins/golang/templates.go
func (t Templater) PluginName() string {
	return "golang/templates"
=======
func (t Validator) PluginName() string {
	return "go/templates/validator"
>>>>>>> tweedy:cli/internal/plugins/golang/templates/validator.go
}
