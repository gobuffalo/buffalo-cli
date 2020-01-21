package built_test

import (
	"github.com/gobuffalo/buffalo-cli/built"
	pop "github.com/gobuffalo/buffalo-cli/cli/internal/plugins/pop/built"
)

var _ built.Initer = &pop.Initer{}
