package soda_test

import (
	"github.com/gobuffalo/buffalo-cli/v2/cli"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/soda"
)

var _ cli.Command = &soda.Cmd{}
