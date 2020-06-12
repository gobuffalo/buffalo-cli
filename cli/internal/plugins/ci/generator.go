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
	f, err := g.buildFiles(root)
	if err != nil {
		return err
	}

	tmpl, err := g.buildTemplate()
	if err != nil {
		return err
	}

	data := struct {
		Name           string
		Database       string
		BuffaloVersion string
	}{
		Name:           name,
		Database:       "postgres",
		BuffaloVersion: "",
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
	defer file.Close()

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

func (g Generator) buildFiles(root string) (*os.File, error) {
	var f *os.File
	var err error

	switch g.provider {
	case "travis":
		f, err = os.Create(filepath.Join(root, ".travis.yml"))
		if err != nil {
			return nil, err
		}
	case "gitlab":
		f, err = os.Create(filepath.Join(root, ".gitlab-ci.yml"))
		if err != nil {
			return nil, err
		}
	case "circleci":
		err = os.MkdirAll(filepath.Join(root, ".circleci"), 0777)
		if err != nil {
			return nil, err
		}

		f, err = os.Create(filepath.Join(root, ".circleci", "config.yml"))
		if err != nil {
			return nil, err
		}
	default:
		err = os.MkdirAll(filepath.Join(root, ".github", "workflows"), 0777)
		if err != nil {
			return nil, err
		}

		f, err = os.Create(filepath.Join(root, ".github", "workflows", "test.yml"))
		if err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (g Generator) templateFile() (pkging.File, error) {
	td := pkger.Include("github.com/gobuffalo/buffalo-cli/v2:/cli/internal/plugins/ci/templates")

	file := g.provider + ".yml.tmpl"
	if g.provider == "" {
		file = "github.yml.tmpl"
	}

	return pkger.Open(filepath.Join(td, file))
}
