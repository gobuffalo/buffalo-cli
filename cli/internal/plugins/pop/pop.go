package pop

import "github.com/gobuffalo/buffalo-cli/internal/plugins"

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&Builder{},
		&Built{},
		&Cmd{},
		&Tester{},
	}
}
