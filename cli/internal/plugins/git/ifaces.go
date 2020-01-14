package git

import (
	"context"
	"os/exec"
)

type VersionRunner interface {
	RunGitVersion(ctx context.Context, cmd *exec.Cmd) (string, error)
}
