package testcmd

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/gobuffalo/buffalo-cli/internal/plugins"
	"github.com/gobuffalo/buffalo-cli/internal/plugins/plugprint"
	"github.com/gobuffalo/here"
)

type TestCmd struct {
	Info      here.Info
	pluginsFn plugins.PluginFeeder
}

var _ plugins.Plugin = &TestCmd{}

func (TestCmd) Name() string {
	return "test"
}

var _ plugprint.Describer = &TestCmd{}

func (TestCmd) Description() string {
	return "Run the tests for the Buffalo app."
}

func (tc *TestCmd) Main(ctx context.Context, args []string) error {
	ioe := plugins.CtxIO(ctx)
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
		return err
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
	return nil
}

func (tc *TestCmd) Cmd(ctx context.Context, args []string) (*exec.Cmd, error) {
	if len(args) == 0 {
		args = append(args, "-cover")
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

func (b *TestCmd) WithHereInfo(i here.Info) {
	b.Info = i
}

func (b *TestCmd) HereInfo() (here.Info, error) {
	if !b.Info.IsZero() {
		return b.Info, nil
	}
	return here.Current()
}

var _ plugins.PluginNeeder = &TestCmd{}

func (b *TestCmd) WithPlugins(f plugins.PluginFeeder) {
	b.pluginsFn = f
}

var _ plugins.PluginScoper = &TestCmd{}

func (bc *TestCmd) ScopedPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	if bc.pluginsFn != nil {
		plugs = bc.pluginsFn()
	}

	var builders []plugins.Plugin
	for _, p := range plugs {
		switch p.(type) {
		case Tester:
			builders = append(builders, p)
		case BeforeTester:
			builders = append(builders, p)
		case AfterTester:
			builders = append(builders, p)
		case Runner:
			builders = append(builders, p)
		case Tagger:
			builders = append(builders, p)
		}
	}
	return builders
}

var _ plugprint.SubCommander = &TestCmd{}

func (bc *TestCmd) SubCommands() []plugins.Plugin {
	var plugs []plugins.Plugin
	for _, p := range bc.ScopedPlugins() {
		if _, ok := p.(Tester); ok {
			plugs = append(plugs, p)
		}
	}
	return plugs
}
