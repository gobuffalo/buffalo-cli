package plugins

// Plugin is the most basic interface a plugin can implement.
type Plugin interface {
	Name() string
}

type Plugins []Plugin

func (p Plugins) ScopedPlugins() []Plugin {
	return []Plugin(p)
}

func (plugs Plugins) ExposedPlugins() []Plugin {
	var exp []Plugin
	for _, p := range plugs {
		if _, ok := p.(Hider); !ok {
			exp = append(exp, p)
		}
	}
	return exp
}

var _ Plugin = Background("")

func Background(name string) Plugin {
	return background(name)
}

type background string

func (b background) Name() string {
	return string(b)
}
