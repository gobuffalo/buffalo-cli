package build

import (
	"context"
	"os/exec"

	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/webpack/internal/ifaces"
	"github.com/gobuffalo/plugins/plugio"
)

// Tooler returns the name of the asset tool to use.
// npm, yarnpkg, etc...
type Tooler = ifaces.Tooler

type Scripter = ifaces.Scripter
type Stdouter = plugio.Outer

type AssetBuilder interface {
	BuildAssets(ctx context.Context, root string, cmd *exec.Cmd) error
}
