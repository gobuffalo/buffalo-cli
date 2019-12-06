package cmdx

import (
	"fmt"
	"io"
	"strings"

	"github.com/gobuffalo/buffalo-cli/cli/plugins"
	"github.com/gobuffalo/buffalo-cli/cli/plugins/plugprint"
	"github.com/spf13/pflag"
)

// Print ...
func Print(w io.Writer, prefix string, main plugins.Plugin, plugs plugins.Plugins, flags *pflag.FlagSet) error {
	p := plugprint.WithUsagePrinter(main, func(w io.Writer) error {
		fmt.Fprintf(w, "Usage of %s:\n", strings.TrimSpace(fmt.Sprintf("%s %s", prefix, main.Name())))
		flags.SetOutput(w)
		flags.PrintDefaults()
		return nil
	})
	return plugprint.Print(w, prefix, p, plugs)
}
