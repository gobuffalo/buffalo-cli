package fizz

import "github.com/gobuffalo/plugins"

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&MigrationGen{},
	}
}
