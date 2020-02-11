package resource

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/resource"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugprint"
	"github.com/gobuffalo/flect/name"
	"github.com/markbates/pkger"
	"github.com/spf13/pflag"
)

var _ plugins.Plugin = &Generator{}
var _ plugprint.FlagPrinter = &Generator{}
var _ resource.Pflagger = &Generator{}
var _ resource.Templater = &Generator{}

type Generator struct {
	modelName string
	flags     *pflag.FlagSet
}

func (g *Generator) PluginName() string {
	return "plush/templates"
}

func (g *Generator) GenerateResourceTemplates(ctx context.Context, root string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must specify a resource name")
	}

	atts, err := attrs.ParseArgs(args[1:]...)
	if err != nil {
		return err
	}

	modelName := g.modelName
	if len(modelName) == 0 {
		modelName = args[0]
	}

	opts := struct {
		Attrs  attrs.Attrs
		Folder string
		Model  name.Ident
		Name   name.Ident
	}{
		Attrs: atts,
		Name:  name.New(args[0]),
		Model: name.New(modelName),
	}
	opts.Folder = opts.Name.Folder().Pluralize().String()

	fp := pkger.Include("github.com/gobuffalo/buffalo-cli/v2:/cli/internal/plugins/plush/generators/resource/templates")

	err = pkger.Walk(fp, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		tf, err := pkger.Open(path)
		if err != nil {
			return err
		}
		defer tf.Close()

		b, err := ioutil.ReadAll(tf)
		if err != nil {
			return err
		}

		tmpl := string(b)

		t, err := template.New(path).Parse(tmpl)
		if err != nil {
			return err
		}

		x := strings.TrimPrefix(path, fp)
		x = strings.TrimSuffix(x, string(filepath.Separator))
		x = filepath.Join(root, "templates", opts.Folder, x)
		x = strings.TrimSuffix(x, ".tmpl")

		if err := os.MkdirAll(filepath.Dir(x), 0755); err != nil {
			return err
		}

		f, err := os.Create(x)
		if err != nil {
			return err
		}
		defer f.Close()

		if err := t.Execute(f, opts); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
