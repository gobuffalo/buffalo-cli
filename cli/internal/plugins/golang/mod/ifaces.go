package mod

import "github.com/gobuffalo/plugins/plugio"

type Stderrer = plugio.Errer
type Stdiner = plugio.Inner
type Stdouter = plugio.Outer

type Requirer interface {
	ModRequire(root string) map[string]string
}

type Replacer interface {
	ModReplacer(root string) map[string]string
}
