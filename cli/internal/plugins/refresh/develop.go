package refresh

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/develop"
	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/gobuffalo/buffalo-cli/v2/plugins/plugprint"
	"github.com/gobuffalo/here"
	"github.com/markbates/refresh/refresh"
	"github.com/spf13/pflag"
)

var _ develop.Developer = &Developer{}
var _ plugins.NamedCommand = &Developer{}
var _ plugins.Plugin = &Developer{}

type Developer struct {
	Debug  bool
	Config string
	help   bool
	flags  *pflag.FlagSet
}

func (dev *Developer) Name() string {
	return "refresh/developer"
}

func (dev *Developer) CmdName() string {
	return "refresh"
}

func (dev *Developer) Develop(ctx context.Context, root string, args []string) error {
	flags := dev.Flags()
	if err := flags.Parse(args); err != nil {
		return err
	}

	if dev.help {
		ioe := plugins.CtxIO(ctx)
		return plugprint.Print(ioe.Stdout(), dev)
	}

	args = flags.Args()

	c, err := dev.config(root)
	if err != nil {
		return err
	}

	info, err := here.Dir(root)
	if err != nil {
		return err
	}

	if len(c.BinaryName) == 0 {
		c.BinaryName = fmt.Sprintf("%s-build", path.Base(info.Module.Path))
	}

	c.Debug = dev.Debug

	r := refresh.NewWithContext(c, ctx)
	r.CommandFlags = args
	return r.Start()
}

func (dev *Developer) config(root string) (*refresh.Configuration, error) {
	if len(dev.Config) == 0 {
		if _, err := os.Stat("./.buffalo.dev.yml"); err == nil {
			dev.Config = "./.buffalo.dev.yml"
		}
	}

	if len(dev.Config) > 0 {
		_, err := os.Stat(dev.Config)
		if err != nil {
			return nil, err
		}
		c := &refresh.Configuration{}
		if err := c.Load(dev.Config); err != nil {
			return nil, err
		}
		return c, nil
	}

	dur, err := time.ParseDuration("200ns")
	if err != nil {
		return nil, err
	}

	c := &refresh.Configuration{
		AppRoot:            root,
		IgnoredFolders:     []string{"vendor", "log", "logs", "webpack", "public", "grifts", "tmp", "bin", "node_modules", ".sass-cache"},
		IncludedExtensions: []string{".go", ".mod", ".env"},
		BuildPath:          "tmp",
		BuildDelay:         dur,
		BinaryName:         "",
		EnableColors:       true,
		LogName:            "buffalo",
	}
	return c, nil
}
