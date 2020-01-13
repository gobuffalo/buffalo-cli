package bzr

import (
	"bytes"
	"context"
	"fmt"
	"strings"
)

// BzrVersioner ...
type BzrVersioner struct {
	versionRunner VersionRunner
}

// Name is the name of the plugin.
// This will also be used for the cli sub-command
// 	"pop" | "heroku" | "auth" | etc...
func (b BzrVersioner) Name() string {
	return "bzr"
}

//Description is a general description of the plugin and its functionalities.
func (b BzrVersioner) Description() string {
	return "Provides bzr related hooks to Buffalo applications."
}

// BuildVersion is used by other commands to get the build
// version of the current source and use it for the build.
func (b *BzrVersioner) BuildVersion(ctx context.Context, root string) (string, error) {
	if b.versionRunner == nil {
		b.versionRunner = &BzrVersionRunner{}
	}

	if ok, err := b.versionRunner.ToolAvailable(); !ok {
		return "", err
	}

	bb := &bytes.Buffer{}
	if err := b.versionRunner.RunVersionCommand(ctx, bb); err != nil {
		return "", fmt.Errorf("%s: %s", err, bb.String())
	}

	s := strings.TrimSpace(bb.String())
	if len(s) == 0 {
		return "", nil
	}
	return s, nil
}
