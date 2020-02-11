package i18n

import (
	"io"

	"github.com/gobuffalo/plugins"
)

type NamedWriter interface {
	plugins.Plugin
	NamedWriter(n string) (io.Writer, error)
}
