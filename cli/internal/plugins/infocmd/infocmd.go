package infocmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo-cli/cli/plugins"
	"github.com/gobuffalo/buffalo-cli/cli/plugins/plugprint"
	"github.com/gobuffalo/buffalo-cli/internal/v1/genny/info"
	"github.com/gobuffalo/clara/genny/rx"
	"github.com/gobuffalo/genny"
	"github.com/spf13/pflag"
)

type InfoCmd struct {
	plugins.IO
	Parent  plugins.Plugin
	Plugins func() plugins.Plugins
	help    bool
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
	for _, p := range ic.informers() {
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

func (ic *InfoCmd) informers() []Informer {
	var plugs []Informer

	if ic.Plugins == nil {
		return nil
	}
	for _, p := range ic.Plugins() {
		if i, ok := p.(Informer); ok {
			plugs = append(plugs, i)
		}
	}

	return plugs
}

// Main implements the `buffalo info` command. Buffalo's checks
// are run first, then any plugins that implement plugins.Informer
// will be run in order at the end.
func (ic *InfoCmd) Main(ctx context.Context, args []string) error {
	out := ic.Stdout()

	flags := pflag.NewFlagSet(ic.String(), pflag.ContinueOnError)
	flags.SetOutput(ioutil.Discard)
	flags.BoolVarP(&ic.help, "help", "h", false, "print this help")
	if err := flags.Parse(args); err != nil {
		return err
	}

	if ic.help {
		ips := ic.informers()
		plugs := make(plugins.Plugins, len(ips))
		for i, ip := range ips {
			plugs[i] = ip
		}
		return plugprint.Print(out, ic, plugs)
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
