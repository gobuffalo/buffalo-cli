package cli

import (
	"bytes"
	"context"
	"testing"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/stretchr/testify/require"
)

func Test_Buffalo_Help(t *testing.T) {
	r := require.New(t)

	b := &Buffalo{}

	ctx := context.Background()

	bb := &bytes.Buffer{}
	ctx = plugins.WithStdout(ctx, bb)

	args := []string{"-h"}

	err := b.Main(ctx, args)
	r.NoError(err)

	r.Contains(bb.String(), b.Description())
}

func Test_Buffalo_Main_SubCommand(t *testing.T) {
	r := require.New(t)

	c := &cp{}
	b := &Buffalo{
		Plugins: plugins.Plugins{
			c,
		},
	}

	ctx := context.Background()

	args := []string{c.Name()}

	exp := []string{"hello"}
	args = append(args, exp...)

	err := b.Main(ctx, args)
	r.NoError(err)
	r.Equal(exp, c.args)
}

func Test_Buffalo_Main_SubCommand_Alias(t *testing.T) {
	r := require.New(t)

	c := &cp{aliases: []string{"sc"}}
	b := &Buffalo{
		Plugins: plugins.Plugins{
			c,
		},
	}

	ctx := context.Background()

	args := []string{"sc"}

	exp := []string{"hello"}
	args = append(args, exp...)

	err := b.Main(ctx, args)
	r.NoError(err)
	r.Equal(exp, c.args)
}
