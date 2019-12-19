package cli

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type cp struct {
	aliases []string
	args    []string
}

func (c *cp) Name() string {
	return "commander"
}

func (c *cp) Main(ctx context.Context, args []string) error {
	c.args = args
	return nil
}

func (c *cp) Aliases() []string {
	return c.aliases
}

func Test_Commands_Find(t *testing.T) {
	r := require.New(t)

	c := &cp{}

	cmds := Commands{c}

	cp, err := cmds.Find(c.Name())
	r.NoError(err)

	exp := []string{"hi"}
	err = cp.Main(context.Background(), exp)
	r.NoError(err)
	r.Equal(exp, c.args)
}

func Test_Commands_Find_Aliases(t *testing.T) {
	r := require.New(t)

	c := &cp{aliases: []string{"hello"}}

	cmds := Commands{c}

	cp, err := cmds.Find(c.Name())
	r.NoError(err)

	exp := []string{"hi"}
	err = cp.Main(context.Background(), exp)
	r.NoError(err)
	r.Equal(exp, c.args)

	cp, err = cmds.Find("hello")
	r.NoError(err)

	exp = []string{"hello", "goodbye"}
	err = cp.Main(context.Background(), exp)
	r.NoError(err)
	r.Equal(exp, c.args)
}
