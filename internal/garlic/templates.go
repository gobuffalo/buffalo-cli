package garlic

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gobuffalo/here"
)

func NewApp(ctx context.Context, root string, args []string) error {
	info, err := here.Dir(root)
	if err != nil {
		return err
	}

	mm := map[string]string{
		filepath.Join(root, "cli", "buffalo.go"):         cliBuffalo,
		filepath.Join(root, "cmd", "buffalo", "main.go"): cliMain,
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

		err = tmpl.Execute(f, struct {
			Name       string
			ImportPath string
		}{
			ImportPath: info.ImportPath,
			Name:       path.Base(info.Module.Path),
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

const cliBuffalo = `
package cli

import (
	"context"
	"fmt"

	"github.com/gobuffalo/buffalo-cli/v2/cli"
	"github.com/gobuffalo/buffalo-cli/v2/plugins"
)

func Buffalo(ctx context.Context, args []string) error {
	fmt.Println("~~~~ Using {{.Name}}/cli.Buffalo ~~~")

	buffalo, err := cli.New()
	if err != nil {
		return err
	}

	buffalo.Plugins = append([]plugins.Plugin{
		// prepend your plugins here
	}, buffalo.Plugins...)

	return buffalo.Main(ctx, args)
}
`

const cliMain = `
package main

import (
	"{{.ImportPath}}/cli"
	"context"
	"log"
	"os"
)

func main() {
	ctx := context.Background()
	if err := cli.Buffalo(ctx, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
`
