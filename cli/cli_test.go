package cli

import (
	"context"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
)

type background string

func (b background) PluginName() string {
	return string(b)
}

var _ plugins.Plugin = &cp{}
var _ plugcmd.Namer = &cp{}
var _ plugcmd.Commander = &cp{}
var _ plugcmd.Aliaser = &cp{}

type cp struct {
	aliases []string
	args    []string
	cmdName string
	name    string
	root    string
}

func (c *cp) PluginName() string {
	if len(c.name) == 0 {
		return "commander"
	}
	return c.name
}

func (c *cp) CmdName() string {
	if len(c.cmdName) == 0 {
		return c.PluginName()
	}
	return c.cmdName
}

func (c *cp) Main(ctx context.Context, root string, args []string) error {
	c.args = args
	c.root = root
	return nil
}

func (c *cp) CmdAliases() []string {
	return c.aliases
}
