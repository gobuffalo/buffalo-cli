package git

import (
	"bytes"
	"context"
	"os/exec"
	"strings"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/build"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugprint"
)

var _ build.Versioner = &Versioner{}
var _ plugins.Plugin = &Versioner{}
var _ plugins.Needer = &Versioner{}
var _ plugins.Scoper = &Versioner{}
var _ plugprint.Describer = &Versioner{}

type Versioner struct {
	pluginsFn plugins.Feeder
}

func (b *Versioner) WithPlugins(f plugins.Feeder) {
	b.pluginsFn = f
}

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

func (b *Versioner) BuildVersion(ctx context.Context, root string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--short", "HEAD")
	bb := &bytes.Buffer{}
	cmd.Stdout = bb

	var fn cmdRunnerFn = func(ctx context.Context, root string, cmd *exec.Cmd) error {
		return cmd.Run()
	}

	for _, p := range b.ScopedPlugins() {
		if vr, ok := p.(CommandRunner); ok {
			fn = vr.RunGitCommand
			break
		}
	}
	if err := fn(ctx, root, cmd); err != nil {
		return "", err
	}
	return strings.TrimSpace(bb.String()), nil
}

func (b Versioner) PluginName() string {
	return "git/versioner"
}

func (b Versioner) Description() string {
	return "Provides the lastest version, SHA, of a Git repo."
}
