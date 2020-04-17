package coreapp

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/newapp"
	"github.com/gobuffalo/here"
	"github.com/gobuffalo/plugins"
	"github.com/markbates/pkger"
)

var _ plugins.Plugin = &Generator{}
var _ newapp.Newapper = &Generator{}

type Generator struct{}

func (g *Generator) PluginName() string {
	return "core/generator"
}

func (g *Generator) Newapp(ctx context.Context, root string, args []string) error {
	td := pkger.Include("github.com/gobuffalo/buffalo-cli/v2:/cli/cmds/newapp/presets/coreapp/_templates")

	her, err := here.Dir(root)
	if err != nil {
		return err
	}

	fmt.Println(">>>TODO cli/cmds/newapp/presets/coreapp/generator.go:36: her ", her)

	err = pkger.Walk(td, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		np := strings.TrimPrefix(path, td)
		np = strings.TrimPrefix(np, "/")
		np = strings.TrimSuffix(np, ".tmpl")
		np = strings.ReplaceAll(np, "-appname-", her.Name)

		if info.IsDir() {
			return os.MkdirAll(filepath.Join(root, np), 0755)
		}

		of, err := pkger.Open(path)
		if err != nil {
			return err
		}
		defer of.Close()

		b, err := ioutil.ReadAll(of)
		if err != nil {
			return err
		}

		f, err := os.Create(filepath.Join(root, np))
		if err != nil {
			return err
		}
		defer f.Close()

		tmpl, err := template.New(path).Parse(string(b))
		if err != nil {
			return err
		}

		if err := tmpl.Execute(f, her); err != nil {
			return err
		}

		io.Copy(f, of)

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
