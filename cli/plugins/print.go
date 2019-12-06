package plugins

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
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
func Print(w io.Writer, prefix string, main Plugin, plugs Plugins) error {
	header := strings.TrimSpace(fmt.Sprintf("%s %s", prefix, main.Name()))
	header = fmt.Sprintf("$ %s", header)
	fmt.Fprintln(w, header)
	for i := 0; i < len(header); i++ {
		fmt.Fprint(w, "-")
	}
	fmt.Fprintln(w)
	if d, ok := main.(Describer); ok {
		fmt.Fprintf(w, "%s\n", d.Description())
	}

	if u, ok := main.(UsagePrinter); ok {
		fmt.Fprintln(w)
		if err := u.PrintUsage(w); err != nil {
			return err
		}
	}

	const ac = "\nAvailable Commands:\n"
	if len(plugs) == 0 {
		return nil
	}
	fmt.Fprint(w, ac)

	sort.Sort(plugs)
	desc := func(p Plugin) string {
		if d, ok := p.(Describer); ok {
			return d.Description()
		}
		return ""
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for _, c := range plugs {
		line := fmt.Sprintf("%s %s %s\t%s\n", prefix, main.Name(), c.Name(), desc(c))
		fmt.Fprintf(tw, "\t%s\n", strings.TrimSpace(line))
	}
	tw.Flush()
	return nil
}
