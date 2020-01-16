package garlic

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/gobuffalo/here"
)

func NewApp(ctx context.Context, root string, args []string) error {
	info, err := here.Dir(root)
	if err != nil {
		return err
	}

	fp := filepath.Join(root, "cli", "buffalo.go")
	if err := os.MkdirAll(filepath.Dir(fp), 0755); err != nil {
		return err
	}

	f, err := os.Create(fp)
	if err != nil {
		return err
	}

	tmpl, err := template.New(fp).Parse(cliBuffalo)
	if err != nil {
		return err
	}

	err = tmpl.Execute(f, struct {
		Name string
	}{
		Name: path.Base(info.Module.Path),
	})

	if err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

const cliBuffalo = `
package cli

import (
	"context"
	"fmt"

	"github.com/gobuffalo/buffalo-cli/cli"
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
