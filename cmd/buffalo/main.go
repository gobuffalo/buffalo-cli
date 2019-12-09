package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/gobuffalo/buffalo-cli/cli"
	"github.com/gobuffalo/here"
	"github.com/markbates/haste"
	"github.com/markbates/jim"
)

func Tidy(ctx context.Context) error {
	c := exec.CommandContext(ctx, "go", "mod", "tidy")
	return c.Run()
}

func main() {
	ctx := context.Background()
	defer Tidy(ctx)

	if err := run(ctx); err != nil {
		Tidy(ctx)
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	info, err := here.Dir(".")
	if err != nil {
		return err
	}

	args := os.Args[1:]

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
