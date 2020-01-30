package cli

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Commands_Find(t *testing.T) {
	r := require.New(t)

	c := &cp{}

	cmds := Commands{c}

	cp, err := cmds.Find(c.PluginName())
	r.NoError(err)

	exp := []string{"hi"}
	err = cp.Main(context.Background(), "", exp)
	r.NoError(err)
	r.Equal(exp, c.args)
}

func Test_Commands_Find_Aliases(t *testing.T) {
	r := require.New(t)

	c := &cp{aliases: []string{"hello"}}

	cmds := Commands{c}

	cp, err := cmds.Find(c.PluginName())
	r.NoError(err)

	exp := []string{"hi"}
	err = cp.Main(context.Background(), "", exp)
	r.NoError(err)
	r.Equal(exp, c.args)

	cp, err = cmds.Find("hello")
	r.NoError(err)

	exp = []string{"hello", "goodbye"}
	err = cp.Main(context.Background(), "", exp)
	r.NoError(err)
	r.Equal(exp, c.args)
}

func Test_Commands_Find_NamedCommand(t *testing.T) {
	r := require.New(t)

	c := &cp{
		cmdName: "hi",
	}

	cmds := Commands{c}

	cp, err := cmds.Find("hi")
	r.NoError(err)

	exp := []string{"hi"}
	err = cp.Main(context.Background(), "", exp)
	r.NoError(err)
	r.Equal(exp, c.args)
}
