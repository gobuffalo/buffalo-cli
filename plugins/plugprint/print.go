package plugprint

import (
	"context"
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
		fmt.Fprintln(w)
		if err := u.PrintFlags(w); err != nil {
			return err
		}
	}

	wp, ok := main.(WithPlugins)
	if !ok {
		return nil
	}

	plugs := wp.WithPlugins()
	if len(plugs) == 0 {
		return nil
	}

	if err := printCommands(w, main, plugs); err != nil {
		return err
	}

	if err := printPlugins(w, main, plugs); err != nil {
		return err
	}

	return nil
}

type WithPlugins interface {
	WithPlugins() plugins.Plugins
}

func printPlugins(w io.Writer, main plugins.Plugin, plugs plugins.Plugins) error {

	wp, ok := main.(WithPlugins)
	if !ok {
		return nil
	}

	plugs = wp.WithPlugins()
	if len(plugs) == 0 {
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "\nUsing Plugins:")
	fmt.Fprintf(tw, "\t%s\t%s\t%s\n", "Type", "Name", "Description")
	fmt.Fprintf(tw, "\t%s\t%s\t%s\n", "----", "----", "-----------")
	for _, p := range plugs {
		fmt.Fprintf(tw, "\t%T\t%s\t%s\n", p, p.Name(), desc(p))
	}

	tw.Flush()
	return nil
}

type Command interface {
	Main(ctx context.Context, args []string) error
}

func printCommands(w io.Writer, main plugins.Plugin, all plugins.Plugins) error {
	if len(all) == 0 {
		return nil
	}

	plugs := make(plugins.Plugins, 0, len(all))
	for _, p := range all {
		if _, ok := p.(Command); ok {
			plugs = append(plugs, p)
		}
	}

	if len(plugs) == 0 {
		return nil
	}

	sort.Sort(plugs)

	const ac = "\nAvailable Commands:\n"
	fmt.Fprint(w, ac)

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for _, c := range plugs {
		fmt.Printf("TYPE: %T\n", c)
		line := fmt.Sprintf("%s\t%s\n", stringer(c), desc(c))
		fmt.Fprintf(tw, "\t%s\n", strings.TrimSpace(line))
	}
	tw.Flush()
	return nil
}

func stringer(p plugins.Plugin) string {
	if st, ok := p.(fmt.Stringer); ok {
		return st.String()
	}
	return p.Name()
}

func desc(p plugins.Plugin) string {
	if d, ok := p.(Describer); ok {
		return d.Description()
	}
	return ""
}
