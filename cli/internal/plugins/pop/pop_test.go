package pop_test

import (
	"github.com/gobuffalo/buffalo-cli/built"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/pop"
)

var _ built.Initer = &pop.Buffalo{}
