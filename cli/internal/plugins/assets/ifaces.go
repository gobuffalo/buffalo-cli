package assets

import (
	"context"
	"os/exec"
)

type AssetBuilder interface {
	BuildAssets(ctx context.Context, cmd *exec.Cmd) error
}
