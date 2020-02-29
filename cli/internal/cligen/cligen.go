package cligen

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/here"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/markbates/safe"
)

type Generator struct {
	Plugins map[string]string
}

func (g *Generator) Generate(ctx context.Context, root string, args []string) error {
	info, err := here.Dir(root)
	if err != nil {
		return err
	}

	x := filepath.Join(root, "cmd", "buffalo")
	mm := map[string]string{
		filepath.Join(x, "cli", "buffalo.go"): tmplBuffalo,
		filepath.Join(x, "main.go"):           tmplMain,
	}

	for fp, body := range mm {
		if err := os.MkdirAll(filepath.Dir(fp), 0755); err != nil {
			return err
		}

		if _, err := os.Stat(fp); err == nil {
			return fmt.Errorf("%s already exists", fp)
		}

		f, err := os.Create(fp)
		if err != nil {
			return err
		}

		body = strings.TrimSpace(body)
		tmpl, err := template.New(fp).Parse(body)
		if err != nil {
			return err
		}

		if g.Plugins == nil {
			g.Plugins = map[string]string{}
		}

		err = tmpl.Execute(f, struct {
			Name       string
			ImportPath string
			Plugs      map[string]string
		}{
			ImportPath: info.Module.Path,
			Name:       path.Base(info.Module.Path),
			Plugs:      g.Plugins,
		})

		if err != nil {
			return err
		}

		if err := f.Close(); err != nil {
			return err
		}
	}

	return g.Run(ctx, root, args)
}

func (g *Generator) Run(ctx context.Context, root string, args []string) error {
	main := filepath.Join(root, "cmd", "newapp")
	if _, err := os.Stat(filepath.Dir(main)); err != nil {
		return err
	}

	bargs := []string{"run", "./cmd/newapp"}
	bargs = append(bargs, args...)

	cmd := exec.CommandContext(ctx, "go", bargs...)
	cmd.Stdin = plugio.Stdin()
	cmd.Stdout = plugio.Stdout()
	cmd.Stderr = plugio.Stderr()
	err := safe.RunE(func() error {
		return cmd.Run()
	})
	if err != nil {
		return err
	}

	return nil
}

const tmplBuffalo = `
package cli

import (
	"context"
	"fmt"

	"github.com/gobuffalo/buffalo-cli/v2/cli"

	{{range $k,$v := .Plugs -}}
	"{{$v}}"
	{{- end}}
)

func Buffalo(ctx context.Context, root string, args []string) error {
	fmt.Print("~~~~ Using {{.Name}}/cmd/buffalo/cli.Buffalo ~~~\n\n")

	buffalo, err := cli.New()
	if err != nil {
		return err
	}

	{{range $k,$v := .Plugs -}}
	buffalo.Plugins = append(buffalo.Plugins, {{$k}}.Plugins()...)
	{{- end}}


	return buffalo.Main(ctx, root, args)
}
`

const tmplMain = `
package main

import (
	"context"
	"log"
	"os"

	"{{.ImportPath}}/cmd/buffalo/cli"
)

func main() {
	ctx := context.Background()
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if err := cli.Buffalo(ctx, pwd, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
`
