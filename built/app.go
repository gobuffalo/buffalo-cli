package built

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gobuffalo/buffalo-cli/cli"
	"github.com/gobuffalo/buffalo-cli/internal/garlic"
	"github.com/gobuffalo/buffalo/runtime"
)

type App struct {
	*cli.Buffalo
	BuildTime    string
	BuildVersion string
	Fallthrough  func(ctx context.Context, args []string) error
	Name         string
	OriginalMain func()
}

func (b *App) Main(ctx context.Context, args []string) error {
	for _, p := range b.Plugins {
		bl, ok := p.(Initer)
		if !ok {
			continue
		}
		if err := bl.BuiltInit(ctx, args); err != nil {
			return err
		}
	}

	if err := b.setBuildInfo(); err != nil {
		return err
	}

	var originalArgs []string
	for i, arg := range args {
		if arg == "--" {
			originalArgs = append([]string{args[0]}, args[i+1:]...)
			args = args[:i]
			break
		}
	}
	if len(args) == 0 {
		if len(originalArgs) != 0 {
			os.Args = originalArgs
		}
		b.OriginalMain()
		return nil
	}

	c := args[0]
	switch c {
	case "version":
		fmt.Printf("%s version %s\n", b.Name, runtime.Build())
		return nil
	}
	if b.Fallthrough != nil {
		return b.Fallthrough(ctx, args)
	}
	return garlic.Run(ctx, args)
}

func (b *App) setBuildInfo() error {
	t, err := time.Parse(time.RFC3339, b.BuildTime)
	if err != nil {
		t = time.Now()
	}
	runtime.SetBuild(runtime.BuildInfo{
		Version: b.BuildVersion,
		Time:    t,
	})
	return nil
}
