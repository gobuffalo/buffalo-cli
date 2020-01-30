package plugins

import (
	"context"
	"fmt"
	"os/exec"
	"path"
)

// Commands is a slice of type `Plugin`
type Commands []Plugin

// Find will try and find the given command in the slice
// by it's `Aliases()`, `CmdName()` or `Name()` methods.
// If it can't be found an error is returned.
func (commands Commands) Find(name string) (Plugin, error) {
	name = path.Base(name)
	for _, c := range commands {
		names := []string{c.PluginName()}
		if a, ok := c.(NamedCommand); ok {
			names = append(names, a.CmdName())
		}
		if a, ok := c.(Aliases); ok {
			names = append(names, a.Aliases()...)
		}
		for _, n := range names {
			if n == name {
				return c, nil
			}
		}
	}
	return nil, fmt.Errorf("command %s not found", name)
}

// Cmd calls the exec.CommandContext, and then sets Stdout, Stderr,
// and Stdin using CtxIO to find IO, if any, in the context to the
// exec.Cmd before returning it.
func Cmd(ctx context.Context, name string, arg ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, name, arg...)

	ioe := CtxIO(ctx)
	cmd.Stdin = ioe.Stdin()
	cmd.Stdout = ioe.Stdout()
	cmd.Stderr = ioe.Stderr()

	return cmd
}
