package git

import (
	"context"
	"fmt"
	"time"

	"github.com/gobuffalo/buffalo-cli/cli/plugins"
)

type Buffalo struct {
	plugins.IO
}

func (b *Buffalo) BuildVersion(ctx context.Context, root string) (string, error) {
	now := time.Now()
	version := now.Format(time.RFC3339)
	fmt.Println(">>>TODO cli/internal/plugins/git/git.go:20: version ", version)
	return "", nil
}

// Name is the name of the plugin.
// This will also be used for the cli sub-command
// 	"pop" | "heroku" | "auth" | etc...
func (b *Buffalo) Name() string {
	return "git"
}
