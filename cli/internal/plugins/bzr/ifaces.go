package bzr

import (
	"context"
	"os/exec"
)

// CommandRunner can be implemented to intercept the
// running of a `bzr` command in any of the plugins in this package.
// This can be useful for testing, logging, wrapping IO, etc...
// It is expected that that plugins in this package will hand over
// control of the exec.Cmd to the first plugin that implements this
// interface.
type Runner interface {
	RunBzr(ctx context.Context, root string, cmd *exec.Cmd) error
}
