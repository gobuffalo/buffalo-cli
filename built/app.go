package built

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gobuffalo/buffalo-cli/internal/garlic"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/buffalo/runtime"
)

type App struct {
	Plugger      plugprint.Plugins
	BuildTime    string
	BuildVersion string
	Fallthrough  func(ctx context.Context, args []string) error
	OriginalMain func()
}

func (b *App) WithPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	if b.Plugger == nil {
		return plugs
	}
	return b.Plugger.WithPlugins()
}

func (b *App) Main(ctx context.Context, args []string) error {
	plugs := b.WithPlugins()

	for _, p := range plugs {
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
		if b.OriginalMain == nil {
			if b.Fallthrough != nil {
				return b.Fallthrough(ctx, args)
			}
			return nil
		}
		if len(originalArgs) != 0 {
			os.Args = originalArgs
		}
		b.OriginalMain()
		return nil
	}

	ioe := plugins.CtxIO(ctx)
	c := args[0]
	switch c {
	case "version":
		fmt.Fprintln(ioe.Stdout(), b.BuildVersion)
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
