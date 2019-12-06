package plugins

import "io"

type UsagePrinter interface {
	PrintUsage(w io.Writer) error
}

type usagePlugin struct {
	Plugin
	fn func(w io.Writer) error
}

func (u usagePlugin) PrintUsage(w io.Writer) error {
	if u.fn == nil {
		return nil
	}
	return u.fn(w)
}

func WithUsagePrinter(p Plugin, fn func(w io.Writer) error) Plugin {
	return usagePlugin{
		Plugin: p,
		fn:     fn,
	}
}
