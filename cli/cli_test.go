package cli

import "context"

type cp struct {
	aliases []string
	args    []string
	cmdName string
	name    string
	root    string
}

func (c *cp) Name() string {
	if len(c.name) == 0 {
		return "commander"
	}
	return c.name
}

func (c *cp) CmdName() string {
	if len(c.cmdName) == 0 {
		return c.Name()
	}
	return c.cmdName
}

func (c *cp) Main(ctx context.Context, root string, args []string) error {
	c.args = args
	c.root = root
	return nil
}

func (c *cp) Aliases() []string {
	return c.aliases
}
