package built

import (
	"context"
	"fmt"
	"os"

	"github.com/gobuffalo/buffalo-cli/v2/internal/garlic"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
)

type App struct {
	Plugger      plugins.Scoper
	BuildTime    string
	BuildVersion string
	Fallthrough  func(ctx context.Context, root string, args []string) error
	OriginalMain func()
}

func (b *App) ScopedPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	if b.Plugger == nil {
		return plugs
	}
	return b.Plugger.ScopedPlugins()
}

func (b *App) Main(ctx context.Context, root string, args []string) error {
	plugs := b.ScopedPlugins()

	for _, p := range plugs {
		bl, ok := p.(Initer)
		if !ok {
			continue
		}
		if err := bl.BuiltInit(ctx, root, args); err != nil {
			return err
		}
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
				return b.Fallthrough(ctx, root, args)
			}
			return nil
		}
		if len(originalArgs) != 0 {
			os.Args = originalArgs
		}
		b.OriginalMain()
		return nil
	}

	c := args[0]
	switch c {
	case "version":
		fmt.Fprintln(plugio.Stdout(plugs...), b.BuildVersion)
		return nil
	}
	if b.Fallthrough != nil {
		return b.Fallthrough(ctx, root, args)
	}
	return garlic.Run(ctx, root, args)
}
