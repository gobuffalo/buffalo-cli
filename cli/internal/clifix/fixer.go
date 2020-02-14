package clifix

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/fix"
	"github.com/gobuffalo/here"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
)

var _ plugins.Plugin = &Fixer{}
var _ plugcmd.Namer = &Fixer{}
var _ fix.Fixer = &Fixer{}

type Fixer struct {
}

func (*Fixer) PluginName() string {
	return "cli/fixer"
}

func (*Fixer) CmdName() string {
	return "cli"
}

func (fixer *Fixer) Fix(ctx context.Context, root string, args []string) error {
	info, err := here.Dir(root)
	if err != nil {
		return err
	}

	x := filepath.Join(root, "cmd", "buffalo")
	mm := map[string]string{
		filepath.Join(x, "cli", "buffalo.go"): cliBuffalo,
		filepath.Join(x, "main.go"):           cliMain,
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
			ImportPath: info.Module.Path,
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
	"github.com/gobuffalo/plugins"
)

func Buffalo(ctx context.Context, root string, args []string) error {
	fmt.Println("~~~~ Using {{.Name}}/cmd/buffalo/cli.Buffalo ~~~\n")

	buffalo, err := cli.New()
	if err != nil {
		return err
	}

	buffalo.Plugins = append([]plugins.Plugin{
		// prepend your plugins here
	}, buffalo.Plugins...)

	return buffalo.Main(ctx, root, args)
}
`

const cliMain = `
package main

import (
	"{{.ImportPath}}/cmd/buffalo/cli"
	"context"
	"log"
	"os"
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
