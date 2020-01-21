package golang

import "github.com/gobuffalo/buffalo-cli/plugins"

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&Templater{},
	}
}
