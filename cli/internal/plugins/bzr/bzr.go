package bzr

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/buildcmd"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
)

const (
	name        = "bzr"
	description = "Provides bzr related hooks to Buffalo applications."
)

//Ensuring bzr is a describer
var _ plugprint.Describer = Buffalo{}

//Ensuring bzr is a buffalo Plugin
var _ plugins.Plugin = Buffalo{}

//Ensuring bzr is a buffalo buildcmd.Versioner
var _ buildcmd.Versioner = &Buffalo{}

//Tconfig is used to mock command things when testing.
var testConfig = struct {
	enabled bool

	resultError   error
	resultVersion string
}{}

type Buffalo struct{}

func (b *Buffalo) runCmd(ctx context.Context, bb *bytes.Buffer) error {
	if testConfig.enabled {
		fmt.Fprint(bb, testConfig.resultVersion)
		return testConfig.resultError
	}

	cmd := exec.CommandContext(ctx, "bzr", "revno")
	cmd.Stdout = bb
	cmd.Stderr = bb

	return cmd.Run()
}

// BuildVersion is used by other commands to get the build
// version of the current source and use it for the build.
func (b *Buffalo) BuildVersion(ctx context.Context, root string) (string, error) {
	if _, err := exec.LookPath("bzr"); err != nil && !testConfig.enabled {
		return "", err
	}

	bb := &bytes.Buffer{}
	if err := b.runCmd(ctx, bb); err != nil {
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
func (b Buffalo) Name() string {
	return name
}

//Description is a general description of the plugin and its functionalities.
func (b Buffalo) Description() string {
	return description
}
