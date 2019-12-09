package fixcmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gobuffalo/buffalo-cli/cli/plugins"
	"github.com/gobuffalo/buffalo-cli/internal/cmdx"
	"github.com/gobuffalo/buffalo-cli/internal/v1/cmd/fix"
)

type FixCmd struct {
	Parent  plugins.Plugin
	Plugins func() plugins.Plugins
	stdin   io.Reader
	stdout  io.Writer
	stderr  io.Writer
}

func (f *FixCmd) SetStderr(w io.Writer) {
	f.stderr = w
}

func (f *FixCmd) SetStdin(r io.Reader) {
	f.stdin = r
}

func (f *FixCmd) SetStdout(w io.Writer) {
	f.stdout = w
}

func (fc *FixCmd) Name() string {
	return "fix"
}

func (fc *FixCmd) Description() string {
	return "Attempt to fix a Buffalo application's API to match version in go.mod"
}

func (f FixCmd) String() string {
	s := f.Name()
	if f.Parent != nil {
		s = fmt.Sprintf("%s %s", f.Parent.Name(), f.Name())
	}
	return strings.TrimSpace(s)
}

// Fix runs any Fixers that are in the Plugins.
// If no arguments are provided it will run all fixers in the Plugins.
// Otherwise Fix will run the fixers for the arguments provided.
// 	buffalo fix
// 	buffalo fix plush pop
// 	buffalo fix -h
func (fc *FixCmd) plugins(ctx context.Context, args []string) error {
	if fc.Plugins == nil {
		return nil
	}
	plugs := fc.Plugins()
	if len(args) > 0 {
		fixers := map[string]Fixer{}
		for _, p := range plugs {
			f, ok := p.(Fixer)
			if !ok {
				continue
			}

			fixers[p.Name()] = f
		}

		for _, a := range args {
			f, ok := fixers[a]
			if !ok {
				return fmt.Errorf("unknown fixer %s", a)
			}
			if err := f.Fix(ctx, []string{}); err != nil {
				return err
			}
		}
		return nil
	}

	for _, p := range plugs {
		f, ok := p.(Fixer)
		if !ok {
			continue
		}

		if err := f.Fix(ctx, args); err != nil {
			return err
		}
	}
	return nil
}

func (fc *FixCmd) Main(ctx context.Context, args []string) error {
	var help bool
	flags := cmdx.NewFlagSet(fc.String())
	flags.BoolVarP(&fix.YesToAll, "yes", "y", false, "update all without asking for confirmation")
	flags.BoolVarP(&help, "help", "h", false, "print this help")

	if err := flags.Parse(args); err != nil {
		return err
	}

	out := fc.stdout
	if out == nil {
		out = os.Stdout
	}

	if help {
		var plugs plugins.Plugins
		if fc.Plugins != nil {
			for _, p := range fc.Plugins() {
				if _, ok := p.(Fixer); ok {
					plugs = append(plugs, p)
				}
			}
		}
		return cmdx.Print(out, fc, plugs, flags)
	}

	if len(args) > 0 {
		return fc.plugins(ctx, args)
	}

	if err := fix.Run(); err != nil {
		return err
	}
	return fc.plugins(ctx, args)
}
