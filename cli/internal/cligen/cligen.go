package cligen

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/here"
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
		filepath.Join(x, "main.go"): tmplMain,
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
	return nil
}

const tmplMain = `
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gobuffalo/buffalo-cli/v2/cli"
)

func main() {
	fmt.Print("~~~~ Using coke/cmd/buffalo ~~~\n\n")

	ctx := context.Background()
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	buffalo, err := cli.New()
	if err != nil {
		log.Fatal(err)
	}

	// append your plugins here
	// buffalo.Plugins = append(buffalo.Plugins, ...)

	err = buffalo.Main(ctx, pwd, os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}
`
