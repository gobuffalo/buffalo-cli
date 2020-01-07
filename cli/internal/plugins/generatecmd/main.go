package generatecmd

import (
	"context"
	"fmt"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/markbates/safe"
)

func (bc *GenerateCmd) beforeGenerate(ctx context.Context, args []string) error {
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

func (bc *GenerateCmd) afterGenerate(ctx context.Context, args []string, err error) error {
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

func (bc *GenerateCmd) Main(ctx context.Context, args []string) error {
	flags := bc.Flags()
	if err := flags.Parse(args); err != nil {
		return err
	}

	ioe := plugins.CtxIO(ctx)
	if len(flags.Args()) == 0 {
		if bc.help {
			return plugprint.Print(ioe.Stdout(), bc)
		}
		return fmt.Errorf("no command provided")
	}

	plugs := bc.ScopedPlugins()

	n := flags.Args()[0]
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
