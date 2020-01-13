package bzr

import (
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/buildcmd"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
)

//Ensuring bzr is a describer
var _ plugprint.Describer = BzrVersioner{}

//Ensuring bzr is a buffalo Plugin
var _ plugins.Plugin = BzrVersioner{}

//Ensuring bzr is a buffalo buildcmd.Versioner
var _ buildcmd.Versioner = &BzrVersioner{}
