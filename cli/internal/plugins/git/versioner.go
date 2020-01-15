package git

import (
	"bytes"
	"context"
	"os/exec"
	"strings"

	"github.com/gobuffalo/buffalo-cli/plugins"
)

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
	if _, err := exec.LookPath("git"); err != nil {
		return "", err
	}

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

// Name is the name of the plugin.
// This will also be used for the cli sub-command
// 	"pop" | "heroku" | "auth" | etc...
func (b Versioner) Name() string {
	return "git/versioner"
}

func (b Versioner) Description() string {
	return "Provides the lastest version, SHA, of a Git repo."
}
