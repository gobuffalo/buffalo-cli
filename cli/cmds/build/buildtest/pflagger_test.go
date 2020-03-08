package buildtest_test

import (
	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/build"
	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/build/buildtest"
)

var _ build.Pflagger = buildtest.Pflagger(nil)
