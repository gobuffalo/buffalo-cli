package garlic

import (
	"context"
	"path"

	"github.com/gobuffalo/buffalo-cli/cli"
	"github.com/gobuffalo/here"
	"github.com/markbates/jim"
)

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
	type tasker interface {
		Task() *jim.Task
	}
	if _, ok := err.(tasker); !ok {
		return err
	}
	b, err := cli.New()
	if err != nil {
		return err
	}
	return b.Main(ctx, args)

}
