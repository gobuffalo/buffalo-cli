package testcmd

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
)

func (tc *TestCmd) Main(ctx context.Context, args []string) error {
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
		return ti.Test(ctx, args[1:])
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
			tc.afterTest(ctx, args, err)
		}
	}()

	if err = tc.beforeTest(ctx, args); err != nil {
		return tc.afterTest(ctx, args, err)
	}

	err = tc.test(ctx, args) // go build ...
	return tc.afterTest(ctx, args, err)

}

func (tc *TestCmd) test(ctx context.Context, args []string) error {
	cmd, err := tc.Cmd(ctx, args)
	if err != nil {
		return err
	}

	for _, p := range tc.ScopedPlugins() {
		if br, ok := p.(Runner); ok {
			return br.RunTests(ctx, cmd)
		}
	}

	return cmd.Run()
}

func (tc *TestCmd) beforeTest(ctx context.Context, args []string) error {
	testers := tc.ScopedPlugins()
	for _, p := range testers {
		if bb, ok := p.(BeforeTester); ok {
			if err := bb.BeforeTest(ctx, args); err != nil {
				return err
			}
		}
	}
	return nil
}

func (tc *TestCmd) afterTest(ctx context.Context, args []string, err error) error {
	testers := tc.ScopedPlugins()
	for _, p := range testers {
		if bb, ok := p.(AfterTester); ok {
			if err := bb.AfterTest(ctx, args, err); err != nil {
				return err
			}
		}
	}
	return err
}

func (tc *TestCmd) Cmd(ctx context.Context, args []string) (*exec.Cmd, error) {
	if len(args) == 0 {
		args = append(args, "./...")
	}

	args, err := tc.buildArgs(ctx, args)
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

func (tc *TestCmd) buildArgs(ctx context.Context, args []string) ([]string, error) {
	args, err := tc.pluginArgs(ctx, args)
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

func (tc *TestCmd) reducePairedArg(key string, args []string) []string {
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

func (tc *TestCmd) pluginArgs(ctx context.Context, args []string) ([]string, error) {
	info, err := tc.HereInfo()
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
