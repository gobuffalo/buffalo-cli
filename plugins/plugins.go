package plugins

// Plugin is the most basic interface a plugin can implement.
type Plugin interface {
	// Name is the name of the plugin.
	// This will also be used for the cli sub-command
	// 	"pop" | "heroku" | "auth" | etc...
	Name() string
}

type Plugins []Plugin

func (p Plugins) ScopedPlugins() []Plugin {
	return []Plugin(p)
}

var _ Plugin = Background("")

func Background(name string) Plugin {
	return background(name)
}

type background string

func (b background) Name() string {
	return string(b)
}

type PluginScoper interface {
	ScopedPlugins() []Plugin
}

type PluginFeeder func() []Plugin

type PluginNeeder interface {
	WithPlugins(PluginFeeder)
}
