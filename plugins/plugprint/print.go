package plugprint

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/gobuffalo/buffalo-cli/plugins"
)

// Print will try and print a helpful Usage printing
// of the plugin and any plugins that are provided.
// 	$ buffalo
// 	---------
//
// 	Usage of buffalo:
// 	-h, --help   print this help
//
// 	Available Commands:
// 		buffalo fix      Attempt to fix a Buffalo application's API to match version in go.mod
// 		buffalo info     Print diagnostic information (useful for debugging)
// 		buffalo version  Print the version information
func Print(w io.Writer, main plugins.Plugin) error {
	if d, ok := main.(Describer); ok {
		fmt.Fprintf(w, "%s\n\n", d.Description())
	}

	header := strings.TrimSpace(fmt.Sprintf("%s", main))
	header = fmt.Sprintf("$ %s", header)
	fmt.Fprintln(w, header)
	for i := 0; i < len(header); i++ {
		fmt.Fprint(w, "-")
	}
	fmt.Fprintln(w)

	if a, ok := main.(Aliases); ok {
		aliases := a.Aliases()
		if len(aliases) != 0 {
			const al = "\nAliases:\n"
			fmt.Fprint(w, al)
			fmt.Fprintln(w, strings.Join(aliases, ", "))
		}
	}

	if u, ok := main.(UsagePrinter); ok {
		fmt.Fprintln(w)
		if err := u.PrintUsage(w); err != nil {
			return err
		}
	}

	if u, ok := main.(FlagPrinter); ok {
		const fp = "\nFlags:\n"
		fmt.Fprint(w, fp)
		if err := u.PrintFlags(w); err != nil {
			return err
		}
	}

	if err := printCommands(w, main); err != nil {
		return err
	}

	if err := printPlugins(w, main); err != nil {
		return err
	}

	return nil
}

func printPlugins(w io.Writer, main plugins.Plugin) error {
	wp, ok := main.(plugins.PluginScoper)
	if !ok {
		return nil
	}

	plugs := wp.ScopedPlugins()
	if len(plugs) == 0 {
		return nil
	}

	fmt.Fprintln(w, "\nUsing Plugins:")
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "\t%s\t%s\t%s\n", "Type", "Name", "Description")
	fmt.Fprintf(tw, "\t%s\t%s\t%s\n", "----", "----", "-----------")
	for _, p := range plugs {
		if _, ok := p.(Hider); ok {
			continue
		}
		fmt.Fprintf(tw, "\t%T\t%s\t%s\n", p, p.Name(), desc(p))
	}

	tw.Flush()
	return nil
}

func printCommands(w io.Writer, main plugins.Plugin) error {
	sc, ok := main.(SubCommander)
	if !ok {
		return nil
	}

	plugs := sc.SubCommands()
	if len(plugs) == 0 {
		return nil
	}

	sort.Slice(plugs, func(i, j int) bool {
		return plugs[i].Name() < plugs[j].Name()
	})

	const ac = "\nAvailable Commands:\n"
	fmt.Fprint(w, ac)

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "\t%s\t%s\n", "Name", "Description")
	fmt.Fprintf(tw, "\t%s\t%s\n", "----", "-----------")
	for _, c := range plugs {
		line := fmt.Sprintf("\t%s\t%s", c.Name(), desc(c))
		fmt.Fprintln(tw, line)
	}
	tw.Flush()
	return nil
}

func desc(p plugins.Plugin) string {
	if d, ok := p.(Describer); ok {
		return d.Description()
	}
	return ""
}
