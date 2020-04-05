package flect

import (
	"github.com/gobuffalo/plugins"
)

const filePath = "/inflections.json"

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&Filer{},
		&Initer{},
	}
}
