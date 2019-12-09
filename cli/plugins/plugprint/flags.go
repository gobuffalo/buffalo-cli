package plugprint

import "io"

type FlagPrinter interface {
	PrintFlags(w io.Writer) error
}
