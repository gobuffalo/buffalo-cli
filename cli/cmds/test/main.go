package test

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/gobuffalo/plugins/plugprint"
)

func (tc *Cmd) Main(ctx context.Context, root string, args []string) error {
	plugs := tc.ScopedPlugins()
	if t := FindTesterFromArgs(args, plugs); t != nil {
		return t.Test(ctx, root, args[1:])
	}

	for _, a := range args {
		if a == "-h" {
			return plugprint.Print(plugio.Stdout(plugs...), tc)
		}
	}

	err := tc.run(ctx, root, args) // go test ...
	return tc.afterTest(ctx, root, args, err)
}

func (tc *Cmd) run(ctx context.Context, root string, args []string) error {
	if err := tc.beforeTest(ctx, root, args); err != nil {
		return plugins.Wrap(tc, err)
	}

	cmd, err := tc.Cmd(ctx, root, args)
	if err != nil {
		return plugins.Wrap(tc, err)
	}

	for _, p := range tc.ScopedPlugins() {
		if br, ok := p.(Runner); ok {
			return br.RunTests(ctx, root, cmd)
		}
	}

	return cmd.Run()
}

func (tc *Cmd) beforeTest(ctx context.Context, root string, args []string) error {
	testers := tc.ScopedPlugins()
	for _, p := range testers {
		if bb, ok := p.(BeforeTester); ok {
			if err := bb.BeforeTest(ctx, root, args); err != nil {
				return plugins.Wrap(bb, err)
			}
		}
	}
	return nil
}

func (tc *Cmd) afterTest(ctx context.Context, root string, args []string, err error) error {
	testers := tc.ScopedPlugins()
	for _, p := range testers {
		if bb, ok := p.(AfterTester); ok {
			if err := bb.AfterTest(ctx, root, args, err); err != nil {
				return plugins.Wrap(bb, err)
			}
		}
	}
	return err
}

func (tc *Cmd) Cmd(ctx context.Context, root string, args []string) (*exec.Cmd, error) {
	if len(args) == 0 {
		args = append(args, "./...")
	}

	args, err := tc.buildArgs(ctx, root, args)
	if err != nil {
		return nil, plugins.Wrap(tc, err)
	}

	cargs := []string{
		"test",
	}
	cargs = append(cargs, args...)

	c := exec.CommandContext(ctx, "go", cargs...)
	fmt.Println(c.Args)

	plugs := tc.ScopedPlugins()
	c.Stdin = plugio.Stdin(plugs...)
	c.Stdout = plugio.Stdout(plugs...)
	c.Stderr = plugio.Stderr(plugs...)
	return c, nil
}

func (tc *Cmd) buildArgs(ctx context.Context, root string, args []string) ([]string, error) {
	args, err := tc.pluginArgs(ctx, root, args)
	if err != nil {
		return nil, plugins.Wrap(tc, err)
	}

	args = tc.reducePairedArg("-tags", args)

	p := args[len(args)-1]

	if strings.HasPrefix(p, ".") {
		return args, nil
	}

	args = append(args, "./...")

	return args, nil
}

func (tc *Cmd) reducePairedArg(key string, args []string) []string {
	nargs := make([]string, 0, len(args))

	ind := -1
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a != key && len(strings.TrimSpace(a)) > 0 {
			nargs = append(nargs, a)
			continue
		}

		if ind == -1 {
			ind = i
			nargs = append(nargs, key, "")
		}

		if len(args) <= i {
			break
		}

		n := args[i+1]
		n = strings.TrimSpace(fmt.Sprintf("%s %s", nargs[ind+1], n))
		nargs[ind+1] = n
		i++
	}
	return nargs
}

func (tc *Cmd) pluginArgs(ctx context.Context, root string, args []string) ([]string, error) {
	plugs := tc.ScopedPlugins()
	for _, p := range plugs {
		bt, ok := p.(Argumenter)
		if !ok {
			continue
		}
		tgs, err := bt.TestArgs(ctx, root)
		if err != nil {
			return nil, plugins.Wrap(bt, err)
		}
		// prepend external build args
		args = append(tgs, args...)
	}
	return args, nil
}
