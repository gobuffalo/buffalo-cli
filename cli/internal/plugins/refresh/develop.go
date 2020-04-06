package refresh

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/develop"
	"github.com/gobuffalo/here"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/gobuffalo/plugins/plugprint"
	"github.com/markbates/refresh/refresh"
	"github.com/spf13/pflag"
)

var _ develop.Developer = &Developer{}
var _ plugcmd.Namer = &Developer{}
var _ plugins.Needer = &Developer{}
var _ plugins.Plugin = &Developer{}
var _ plugins.Scoper = &Developer{}

type Developer struct {
	config    string
	debug     bool
	flags     *pflag.FlagSet
	help      bool
	pluginsFn plugins.Feeder
}

func (dev *Developer) PluginName() string {
	return "refresh/developer"
}

func (dev *Developer) CmdName() string {
	return "refresh"
}

func (dev *Developer) WithPlugins(f plugins.Feeder) {
	dev.pluginsFn = f
}

func (dev *Developer) ScopedPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin

	if dev.pluginsFn == nil {
		return plugs
	}

	for _, p := range dev.pluginsFn() {
		switch p.(type) {
		case Tagger:
			plugs = append(plugs, p)
		case Stdouter:
			plugs = append(plugs, p)
		}
	}
	return plugs
}

func (dev *Developer) Develop(ctx context.Context, root string, args []string) error {
	flags := dev.Flags()
	if err := flags.Parse(args); err != nil {
		return plugins.Wrap(dev, err)
	}

	if dev.help {
		return plugprint.Print(plugio.Stdout(dev.ScopedPlugins()...), dev)
	}

	args = flags.Args()

	c, err := dev.buildConfig(ctx, root)
	if err != nil {
		return plugins.Wrap(dev, err)
	}

	info, err := here.Dir(root)
	if err != nil {
		return plugins.Wrap(dev, err)
	}

	if len(c.BinaryName) == 0 {
		c.BinaryName = fmt.Sprintf("%s-build", path.Base(info.Module.Path))
	}

	c.Debug = dev.debug

	r := refresh.NewWithContext(c, ctx)
	r.CommandFlags = args
	return r.Start()
}

func (dev *Developer) buildConfig(ctx context.Context, root string) (*refresh.Configuration, error) {
	if len(dev.config) == 0 {
		if _, err := os.Stat("./.buffalo.dev.yml"); err == nil {
			dev.config = "./.buffalo.dev.yml"
		}
	}

	if len(dev.config) > 0 {
		_, err := os.Stat(dev.config)
		if err != nil {
			return nil, plugins.Wrap(dev, err)
		}
		c := &refresh.Configuration{}
		if err := c.Load(dev.config); err != nil {
			return nil, plugins.Wrap(dev, err)
		}
		return c, nil
	}

	dur, err := time.ParseDuration("200ns")
	if err != nil {
		return nil, plugins.Wrap(dev, err)
	}

	var bflags []string
	tags, err := dev.buildTags(ctx, root)
	if err != nil {
		return nil, plugins.Wrap(dev, err)
	}

	if len(tags) > 0 {
		bflags = append(bflags, "-tags", strings.Join(tags, " "))
	}

	c := &refresh.Configuration{
		AppRoot:            root,
		IgnoredFolders:     []string{"vendor", "log", "logs", "webpack", "public", "grifts", "tmp", "bin", "node_modules", ".sass-cache"},
		IncludedExtensions: []string{".go", ".mod", ".env"},
		BuildPath:          "tmp",
		BuildDelay:         dur,
		BuildFlags:         bflags,
		BinaryName:         "",
		EnableColors:       true,
		LogName:            "buffalo",
	}
	return c, nil
}

func (dev *Developer) buildTags(ctx context.Context, root string) ([]string, error) {
	var tags []string
	for _, p := range dev.ScopedPlugins() {
		t, ok := p.(Tagger)
		if !ok {
			continue
		}
		bt, err := t.BuildTags(ctx, root)
		if err != nil {
			return nil, plugins.Wrap(dev, err)
		}
		tags = append(tags, bt...)
	}

	return tags, nil
}
