package plugins

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
)

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

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for _, c := range plugs {
		line := fmt.Sprintf("%s %s %s\t%s\n", prefix, main.Name(), c.Name(), Description(c))
		fmt.Fprintf(tw, "\t%s\n", strings.TrimSpace(line))
	}
	tw.Flush()
	return nil
}
