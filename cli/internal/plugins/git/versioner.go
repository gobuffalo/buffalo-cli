package git

import (
	"bytes"
	"context"
	"os/exec"
	"strings"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/build"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
)

var _ build.Versioner = &Versioner{}
var _ plugins.Plugin = &Versioner{}
var _ plugins.PluginNeeder = &Versioner{}
var _ plugins.PluginScoper = &Versioner{}
var _ plugprint.Describer = &Versioner{}

type Versioner struct {
	pluginsFn plugins.PluginFeeder
}

func (b *Versioner) WithPlugins(f plugins.PluginFeeder) {
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

	var fn cmdRunnerFn = func(ctx context.Context, cmd *exec.Cmd) error {
		return cmd.Run()
	}

	for _, p := range b.ScopedPlugins() {
		if vr, ok := p.(CommandRunner); ok {
			fn = vr.RunGitCommand
			break
		}
	}
	if err := fn(ctx, cmd); err != nil {
		return "", err
	}
	return strings.TrimSpace(bb.String()), nil
}

func (b Versioner) Name() string {
	return "git/versioner"
}

func (b Versioner) Description() string {
	return "Provides the lastest version, SHA, of a Git repo."
}
