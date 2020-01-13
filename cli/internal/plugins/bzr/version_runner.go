package bzr

import (
	"bytes"
	"context"
)

//VersionRunner interface allows to swap the way we determine version command availability and output.
type VersionRunner interface {
	ToolAvailable() (bool, error)
	RunVersionCommand(context.Context, *bytes.Buffer) error
}
