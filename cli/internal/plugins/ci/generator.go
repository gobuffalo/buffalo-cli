package ci

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/markbates/pkger"
	"github.com/markbates/pkger/pkging"
	"github.com/spf13/pflag"
)

type Generator struct {
	provider string
	flags    *pflag.FlagSet
}

func (Generator) PluginName() string {
	return "ci"
}

func (Generator) Description() string {
	return "Generates CI configuration file"
}

func (g Generator) Newapp(ctx context.Context, root string, name string, args []string) error {
	err := os.MkdirAll(filepath.Join(root, ".github", "workflows"), 0777)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(root, ".github", "workflows", "test.yml"))
	if err != nil {
		return err
	}
	defer f.Close()

	tmpl, err := g.buildTemplate()
	if err != nil {
		return err
	}

	data := struct {
		Name string
	}{
		Name: name,
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		return err
	}

	return nil
}

func (g Generator) buildTemplate() (*template.Template, error) {
	file, err := g.templateFile()
	if err != nil {
		return nil, err
	}

	t, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	template, err := template.New("ci-file").Parse(string(t))
	if err != nil {
		return nil, err
	}

	return template, nil
}

func (g Generator) templateFile() (pkging.File, error) {
	td := pkger.Include("github.com/gobuffalo/buffalo-cli/v2:/cli/internal/plugins/ci/templates")

	file := "github.yml.tmpl"
	// if g.style == "standard" {
	// 	file = "Dockerfile.standard"
	// }

	return pkger.Open(filepath.Join(td, file))
}
