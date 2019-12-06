package plugprint

import (
	"io"

	"github.com/gobuffalo/buffalo-cli/cli/plugins"
)

// UsagePrinter is called by `Print` and can be implemented
// to print a large block of usage information after the
// `Describer` interface is called. This is useful for printing
// flag information, links, and other messages to users.
type UsagePrinter interface {
	PrintUsage(w io.Writer) error
}

type usagePlugin struct {
	plugins.Plugin
	fn func(w io.Writer) error
}

func (u usagePlugin) PrintUsage(w io.Writer) error {
	if u.fn == nil {
		return nil
	}
	return u.fn(w)
}

// WithUsagePrinter wraps the provided Plugin with a plugin
// that implements the UsagePrinter interface using the provided
// function to fill the interface.
func WithUsagePrinter(p plugins.Plugin, fn func(w io.Writer) error) plugins.Plugin {
	return usagePlugin{
		Plugin: p,
		fn:     fn,
	}
}
