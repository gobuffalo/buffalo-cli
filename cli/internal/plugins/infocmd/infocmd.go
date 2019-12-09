package infocmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo-cli/cli/plugins"
	"github.com/gobuffalo/buffalo-cli/internal/cmdx"
	"github.com/gobuffalo/buffalo-cli/internal/v1/genny/info"
	"github.com/gobuffalo/clara/genny/rx"
	"github.com/gobuffalo/genny"
)

type InfoCmd struct {
	Parent  plugins.Plugin
	Plugins func() plugins.Plugins
	stdin   io.Reader
	stdout  io.Writer
	stderr  io.Writer
	help    bool
}

func (i *InfoCmd) SetStderr(w io.Writer) {
	i.stderr = w
}

func (i *InfoCmd) SetStdin(r io.Reader) {
	i.stdin = r
}

func (i *InfoCmd) SetStdout(w io.Writer) {
	i.stdout = w
}

func (ic *InfoCmd) Name() string {
	return "info"
}

func (ic *InfoCmd) Description() string {
	return "Print diagnostic information (useful for debugging)"
}

func (i InfoCmd) String() string {
	s := i.Name()
	if i.Parent != nil {
		s = fmt.Sprintf("%s %s", i.Parent.Name(), i.Name())
	}
	return strings.TrimSpace(s)
}

// Info runs all of the plugins that implement the
// `Informer` interface in order.
func (ic *InfoCmd) plugins(ctx context.Context, args []string) error {
	if ic.Plugins == nil {
		return nil
	}
	plugs := ic.Plugins()
	for _, p := range plugs {
		i, ok := p.(Informer)
		if !ok {
			continue
		}
		if err := i.Info(ctx, args); err != nil {
			return err
		}
	}
	return nil
}

// Main implements the `buffalo info` command. Buffalo's checks
// are run first, then any plugins that implement plugins.Informer
// will be run in order at the end.
func (ic *InfoCmd) Main(ctx context.Context, args []string) error {
	out := ic.stdout
	if out == nil {
		out = os.Stdout
	}

	flags := cmdx.NewFlagSet(ic.String())
	flags.BoolVarP(&ic.help, "help", "h", false, "print this help")
	if err := flags.Parse(args); err != nil {
		return err
	}

	if ic.help {
		return cmdx.Print(out, ic, nil, flags)
	}

	args = flags.Args()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	run := genny.WetRunner(ctx)

	opts := &rx.Options{
		Out: rx.NewWriter(out),
	}
	if err := run.WithNew(rx.New(opts)); err != nil {
		return err
	}

	iopts := &info.Options{
		Out: rx.NewWriter(out),
	}

	if err := run.WithNew(info.New(iopts)); err != nil {
		return err
	}

	if err := run.Run(); err != nil {
		return err
	}
	return ic.plugins(ctx, args)
}
