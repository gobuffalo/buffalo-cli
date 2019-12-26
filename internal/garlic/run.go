package garlic

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/gobuffalo/buffalo-cli/cli"
	"github.com/gobuffalo/here"
	"github.com/markbates/jim"
)

type tasker interface {
	Task() *jim.Task
}

func Run(ctx context.Context, args []string) error {
	info, err := here.Dir(".")
	if err != nil {
		return err
	}

	ip := path.Join(info.Module.Path, "cli")
	t := &jim.Task{
		Info: info,
		Args: args,
		Pkg:  ip,
		Sel:  "cli",
		Name: "Buffalo",
	}

	err = jim.Run(ctx, t)
	if err == nil {
		return nil
	}

	if _, ok := err.(tasker); !ok {
		return err
	}

	fp := filepath.Join(info.Root, "cli", "buffalo.go")
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

	err = jim.Run(ctx, t)
	if err == nil {
		return nil
	}

	b, err := cli.New()
	if err != nil {
		return err
	}
	return b.Main(ctx, args)

}
