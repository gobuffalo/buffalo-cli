package generatecmd

import (
	"context"
	"fmt"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/markbates/safe"
)

func (bc *Command) beforeGenerate(ctx context.Context, args []string) error {
	builders := bc.ScopedPlugins()
	for _, p := range builders {
		if bb, ok := p.(BeforeGenerator); ok {
			err := safe.RunE(func() error {
				return bb.BeforeGenerate(ctx, args)
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (bc *Command) afterGenerate(ctx context.Context, args []string, err error) error {
	builders := bc.ScopedPlugins()
	for _, p := range builders {
		if bb, ok := p.(AfterGenerator); ok {
			err := safe.RunE(func() error {
				return bb.AfterGenerate(ctx, args, err)
			})
			if err != nil {
				return err
			}
		}
	}
	return err
}

// Main implements cli.Command and is the entry point for `buffalo generate`
func (bc *Command) Main(ctx context.Context, args []string) error {
	ioe := plugins.CtxIO(ctx)
	if len(args) == 0 {
		if err := plugprint.Print(ioe.Stdout(), bc); err != nil {
			return err
		}
		return fmt.Errorf("no command provided")
	}

	if len(args) == 1 && args[0] == "-h" {
		return plugprint.Print(ioe.Stdout(), bc)
	}

	plugs := bc.ScopedPlugins()

	n := args[0]
	cmds := plugins.Commands(plugs)
	p, err := cmds.Find(n)
	if err != nil {
		return err
	}

	b, ok := p.(Generator)
	if !ok {
		return fmt.Errorf("unknown command %q", n)
	}

	if err = bc.beforeGenerate(ctx, args); err != nil {
		return bc.afterGenerate(ctx, args, err)
	}

	err = safe.RunE(func() error {
		return b.Generate(ctx, args[1:])
	})

	return bc.afterGenerate(ctx, args, err)
}
