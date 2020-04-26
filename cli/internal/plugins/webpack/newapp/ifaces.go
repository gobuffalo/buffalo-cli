package newapp

import (
	"context"
	"os/exec"

	"github.com/gobuffalo/plugins/plugio"
)

type Stdouter = plugio.Outer
type Stdiner = plugio.Inner
type Stderrer = plugio.Errer

// CommandRunner can be implemented to intercept the
// running of a `webpack` command in any of the plugins in this package.
// This can be useful for testing, logging, wrapping IO, etc...
// It is expected that that plugins in this package will hand over
// control of the exec.Cmd to the first plugin that implements this
// interface.
type CommandRunner interface {
	RunWebpackCommand(ctx context.Context, root string, cmd *exec.Cmd) error
}
