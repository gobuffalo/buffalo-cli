package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"os/signal"
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

	// trap Ctrl+C and call cancel on the context
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

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
