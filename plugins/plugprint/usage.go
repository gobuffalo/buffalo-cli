package plugprint

import (
	"io"
)

// UsagePrinter is called by `Print` and can be implemented
// to print a large block of usage information after the
// `Describer` interface is called. This is useful for printing
// flag information, links, and other messages to users.
type UsagePrinter interface {
	PrintUsage(w io.Writer) error
}
