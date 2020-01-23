package test

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/gobuffalo/buffalo-cli/v2/plugins/plugprint"
	"github.com/gobuffalo/here"
)

func (tc *Cmd) Main(ctx context.Context, root string, args []string) error {
	ioe := plugins.CtxIO(ctx)

	plugs := tc.ScopedPlugins()

	var ti Tester
	if len(args) > 0 {
		n := args[0]
		cmds := plugins.Commands(plugs)
		p, err := cmds.Find(n)
		if err == nil {
			var ok bool
			ti, ok = p.(Tester)
			if !ok {
				return fmt.Errorf("unknown command %q", n)
			}
		}
	}
	if ti != nil {
		return ti.Test(ctx, root, args[1:])
	}

	for _, a := range args {
		if a == "-h" {
			return plugprint.Print(ioe.Stdout(), tc)
		}
	}

	var err error
	defer func() {
		if e := recover(); e != nil {
			var ok bool
			err, ok = e.(error)
			if !ok {
				err = fmt.Errorf("%s", e)
			}
			tc.afterTest(ctx, root, args, err)
		}
	}()

	if err = tc.beforeTest(ctx, root, args); err != nil {
		return tc.afterTest(ctx, root, args, err)
	}

	err = tc.test(ctx, root, args) // go build ...
	return tc.afterTest(ctx, root, args, err)

}

func (tc *Cmd) test(ctx context.Context, root string, args []string) error {
	cmd, err := tc.Cmd(ctx, root, args)
	if err != nil {
		return err
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
				return err
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
				return err
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
		return nil, err
	}

	cargs := []string{
		"test",
	}
	cargs = append(cargs, args...)

	c := exec.CommandContext(ctx, "go", cargs...)
	fmt.Println(c.Args)

	ioe := plugins.CtxIO(ctx)
	c.Stdin = ioe.Stdin()
	c.Stdout = ioe.Stdout()
	c.Stderr = ioe.Stderr()
	return c, nil
}

func (tc *Cmd) buildArgs(ctx context.Context, root string, args []string) ([]string, error) {
	args, err := tc.pluginArgs(ctx, root, args)
	if err != nil {
		return nil, err
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
	info, err := here.Dir(root)
	if err != nil {
		return nil, err
	}

	plugs := tc.ScopedPlugins()
	for _, p := range plugs {
		bt, ok := p.(Argumenter)
		if !ok {
			continue
		}
		tgs, err := bt.TestArgs(ctx, info.Dir)
		if err != nil {
			return nil, err
		}
		// prepend external build args
		args = append(tgs, args...)
	}
	return args, nil
}
