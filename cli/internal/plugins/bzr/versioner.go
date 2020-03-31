package bzr

import (
	"bytes"
	"context"
	"os/exec"
	"strings"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/build"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/gobuffalo/plugins/plugprint"
)

var _ build.Versioner = &Versioner{}
var _ plugins.Plugin = Versioner{}
var _ plugins.Needer = &Versioner{}
var _ plugprint.Describer = Versioner{}

// Versioner ...
type Versioner struct {
	pluginsFn plugins.Feeder
}

// WithPlugins ...
func (b *Versioner) WithPlugins(f plugins.Feeder) {
	b.pluginsFn = f
}

// ScopedPlugins ...
func (b *Versioner) ScopedPlugins() []plugins.Plugin {
	if b.pluginsFn == nil {
		return []plugins.Plugin{}
	}

	var scoped []plugins.Plugin
	for _, p := range b.pluginsFn() {
		switch p.(type) {
		case Runner:
			scoped = append(scoped, p)
		}
	}

	return scoped
}

// BuildVersion is used by other commands to get the build
// version of the current source and use it for the build.
func (b *Versioner) BuildVersion(ctx context.Context, root string) (string, error) {
	plugs := b.ScopedPlugins()

	cmd := exec.CommandContext(ctx, "bzr", "revno")

	bb := &bytes.Buffer{}
	cmd.Stdout = bb
	cmd.Stderr = plugio.Stderr(plugs...)

	fn := func(ctx context.Context, root string, cmd *exec.Cmd) error {
		return cmd.Run()
	}

	for _, p := range plugs {
		if vr, ok := p.(Runner); ok {
			fn = vr.RunBzr
			break
		}
	}

	if err := fn(ctx, root, cmd); err != nil {
		return "", err
	}

	return strings.TrimSpace(bb.String()), nil
}

// Name is the name of the plugin.
// This will also be used for the cli sub-command
// 	"pop" | "heroku" | "auth" | etc...
func (b Versioner) PluginName() string {
	return "bzr"
}

//Description is a general description of the plugin and its functionalities.
func (b Versioner) Description() string {
	return "Provides bzr related hooks to Buffalo applications."
}
