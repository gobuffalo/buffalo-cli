package generatecmd

import (
	"context"
	"fmt"

	"github.com/gobuffalo/buffalo-cli/internal/plugins"
	"github.com/gobuffalo/buffalo-cli/internal/plugins/plugprint"
)

func (bc *GenerateCmd) beforeGenerate(ctx context.Context, args []string) error {
	builders := bc.ScopedPlugins()
	for _, p := range builders {
		if bb, ok := p.(BeforeGenerator); ok {
			if err := bb.BeforeGenerate(ctx, args); err != nil {
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
			if err := bb.AfterGenerate(ctx, args, err); err != nil {
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

	plugs := bc.ScopedPlugins()

	if len(flags.Args()) > 0 {
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
		return b.Generate(ctx, args[1:])
	}

	if bc.help {
		return plugprint.Print(ioe.Stdout(), bc)
	}

	var err error
	defer func() {
		if e := recover(); e != nil {
			var ok bool
			err, ok = e.(error)
			if !ok {
				err = fmt.Errorf("%s", e)
			}
			bc.afterGenerate(ctx, args, err)
		}
	}()

	if err = bc.beforeGenerate(ctx, args); err != nil {
		return bc.afterGenerate(ctx, args, err)
	}

	// err = bc.build(ctx) // go build ...
	return bc.afterGenerate(ctx, args, err)
}
