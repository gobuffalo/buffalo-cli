package garlic

import (
	"context"
	"path"

	"github.com/gobuffalo/buffalo-cli/cli"
	"github.com/gobuffalo/here"
	"github.com/markbates/haste"
	"github.com/markbates/jim"
)

func Run(ctx context.Context, args []string) error {
	info, err := here.Dir(".")
	if err != nil {
		return err
	}

	ip := path.Join(info.Module.Path, "cli")
	h, err := haste.New(ip)
	if err != nil {
		return err
	}

	const bufFn = "func Buffalo(context.Context, []string) error"

	if _, err := h.Funcs().Find(bufFn); err != nil {
		b, err := cli.New()
		if err != nil {
			return err
		}
		return b.Main(ctx, args)
	}

	t := &jim.Task{
		Info: info,
		Args: args,
		Pkg:  ip,
		Sel:  "cli",
		Name: "Buffalo",
	}

	return jim.Run(ctx, t)
}
