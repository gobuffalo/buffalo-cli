package bzr

import (
	"bytes"
	"context"
	"os/exec"
	"strings"

	"github.com/gobuffalo/plugins"
)

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
		case CommandRunner:
			scoped = append(scoped, p)
		}
	}

	return scoped
}

// BuildVersion is used by other commands to get the build
// version of the current source and use it for the build.
func (b *Versioner) BuildVersion(ctx context.Context, root string) (string, error) {

	cmd := exec.CommandContext(ctx, "bzr", "revno")
	bb := &bytes.Buffer{}
	cmd.Stdout = bb

	var fn cmdRunnerFn = func(ctx context.Context, root string, cmd *exec.Cmd) error {
		return cmd.Run()
	}

	for _, p := range b.ScopedPlugins() {
		if vr, ok := p.(CommandRunner); ok {
			fn = vr.RunBzrCommand
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
