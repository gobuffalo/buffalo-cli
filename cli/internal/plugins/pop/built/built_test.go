package built_test

import (
	"github.com/gobuffalo/buffalo-cli/v2/built"
	pop "github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/pop/built"
)

var _ built.Initer = &pop.Initer{}
