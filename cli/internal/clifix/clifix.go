package clifix

import "github.com/gobuffalo/plugins"

//Plugins returns the plugins for the clifix.
func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&Fixer{},
	}
}
