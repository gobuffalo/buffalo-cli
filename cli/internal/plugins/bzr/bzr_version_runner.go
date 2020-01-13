package bzr

import (
	"bytes"
	"context"
	"os/exec"
)

var _ VersionRunner = &BzrVersionRunner{}

// BzrVersionRunner ...
type BzrVersionRunner struct{}

// ToolAvailable checks if the bzr exec is available
func (b *BzrVersionRunner) ToolAvailable() (bool, error) {
	if _, err := exec.LookPath("bzr"); err != nil {
		return false, err
	}

	return true, nil
}

// RunVersionCommand runs bzr command to later on extract the version from the passed buffer.
func (b *BzrVersionRunner) RunVersionCommand(ctx context.Context, bb *bytes.Buffer) error {
	cmd := exec.CommandContext(ctx, "bzr", "revno")
	cmd.Stdout = bb
	cmd.Stderr = bb

	return cmd.Run()
}
