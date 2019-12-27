package plugins

import "io"

type PluginScoper interface {
	ScopedPlugins() []Plugin
}

type PluginFeeder func() []Plugin

type PluginNeeder interface {
	WithPlugins(PluginFeeder)
}

type Hider interface {
	HidePlugin()
}

type IO interface {
	StderrGetter
	StdinGetter
	StdoutGetter
}

type StdinGetter interface {
	Stdin() io.Reader
}

type StdoutGetter interface {
	Stdout() io.Writer
}

type StderrGetter interface {
	Stderr() io.Writer
}
