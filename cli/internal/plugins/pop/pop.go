package pop

import "github.com/gobuffalo/buffalo-cli/plugins"

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&Builder{},
		&Built{},
		&Cmd{},
		&ModelGen{},
		&Tester{},
	}
}
