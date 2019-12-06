package cmdx

import (
	"io"

	"github.com/gobuffalo/buffalo-cli/cli/plugins"
	"github.com/spf13/pflag"
)

func Print(w io.Writer, prefix string, main plugins.Plugin, plugs plugins.Plugins, flags *pflag.FlagSet) error {
	p := plugins.WithUsagePrinter(main, func(w io.Writer) error {
		flags.SetOutput(w)
		flags.Usage()
		return nil
	})
	return plugins.Print(w, prefix, p, plugs)
}
