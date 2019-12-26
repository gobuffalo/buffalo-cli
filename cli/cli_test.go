package cli

import "context"

type cp struct {
	name    string
	cmdName string
	aliases []string
	args    []string
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

func (c *cp) Main(ctx context.Context, args []string) error {
	c.args = args
	return nil
}

func (c *cp) Aliases() []string {
	return c.aliases
}
