package bzr

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type Buffalo struct {
}

func (b *Buffalo) BuildVersion(ctx context.Context, root string) (string, error) {
	if _, err := exec.LookPath("bzr"); err != nil {
		return "", err
	}

	bb := &bytes.Buffer{}

	cmd := exec.CommandContext(ctx, "bzr", "revno")
	cmd.Stdout = bb
	cmd.Stderr = bb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s: %s", err, bb.String())
	}
	s := strings.TrimSpace(bb.String())
	if len(s) == 0 {
		return "", nil
	}
	return s, nil
}

// Name is the name of the plugin.
// This will also be used for the cli sub-command
// 	"pop" | "heroku" | "auth" | etc...
func (b *Buffalo) Name() string {
	return "bzr"
}

func (b *Buffalo) Description() string {
	return "Provides bzr related hooks to Buffalo applications."
}
