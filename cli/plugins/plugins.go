package plugins

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
)

// Plugin is the most basic interface a plugin can implement.
type Plugin interface {
	// Name is the name of the plugin.
	// This will also be used for the cli sub-command
	// 	"pop" | "heroku" | "auth" | etc...
	Name() string
}

type Plugins []Plugin

// Len is the number of elements in the collection.
func (plugs Plugins) Len() int {
	return len(plugs)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (plugs Plugins) Less(i int, j int) bool {
	return plugs[i].Name() < plugs[j].Name()
}

// Swap swaps the elements with indexes i and j.
func (plugs Plugins) Swap(i int, j int) {
	plugs[i], plugs[j] = plugs[j], plugs[i]
}

func Print(w io.Writer, prefix string, main Plugin, plugs Plugins) error {
	header := strings.TrimSpace(fmt.Sprintf("%s %s", prefix, main.Name()))
	fmt.Fprintf(w, "%s\n", header)
	if d, ok := main.(Describer); ok {
		fmt.Fprintf(w, "%s\n", d.Description())
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
