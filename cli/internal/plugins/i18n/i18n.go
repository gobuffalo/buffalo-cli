package i18n

import (
	"github.com/gobuffalo/plugins"
)

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&Generator{},
		&Newapp{},
	}
}
