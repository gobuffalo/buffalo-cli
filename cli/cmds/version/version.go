package version

import "github.com/gobuffalo/plugins"

var Version = "buffalo-cli/unknown"

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&Cmd{},
	}
}
