package git

import (
	"bytes"
	"context"
	"os/exec"
	"strings"

	"github.com/gobuffalo/buffalo-cli/plugins"
)

type Buffalo struct {
	plugins.IO
}

func (b *Buffalo) BuildVersion(ctx context.Context, root string) (string, error) {
	if _, err := exec.LookPath("git"); err != nil {
		return "", err
	}

	bb := &bytes.Buffer{}

	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--short", "HEAD")
	cmd.Stdout = bb
	if err := cmd.Run(); err != nil {
		return "", err
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
	return "git"
}

func (b *Buffalo) Description() string {
	return "Provides git related hooks to Buffalo applications."
}
