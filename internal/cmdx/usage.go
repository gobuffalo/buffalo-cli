package cmdx

import (
	"io"

	"github.com/gobuffalo/buffalo-cli/cli/plugins"
	"github.com/gobuffalo/buffalo-cli/cli/plugins/plugprint"
)

// Print ...
func Print(w io.Writer, main plugins.Plugin, plugs plugins.Plugins, flags *FlagSet) error {
	// p := plugprint.WithUsagePrinter(main, func(w io.Writer) error {
	// 	fmt.Fprintf(w, "Usage of %s:\n", strings.TrimSpace(fmt.Sprintf("%s %s", prefix, main.Name())))
	// 	flags.SetOutput(w)
	// 	flags.PrintDefaults()
	// 	return nil
	// })
	return plugprint.Print(w, main, plugs)
}
