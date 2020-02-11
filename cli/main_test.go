package cli

import (
	"bytes"
	"context"
	"testing"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/stretchr/testify/require"
)

func Test_Buffalo_Help(t *testing.T) {
	r := require.New(t)

	stdout := &bytes.Buffer{}

	b := &Buffalo{
		Plugins: plugins.Plugins{
			plugio.NewOuter(stdout),
		},
	}

	ctx := context.Background()

	args := []string{"-h"}

	err := b.Main(ctx, "", args)
	r.NoError(err)

	r.Contains(stdout.String(), b.Description())
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

	args := []string{c.PluginName()}

	exp := []string{"hello"}
	args = append(args, exp...)

	err := b.Main(ctx, "", args)
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

	err := b.Main(ctx, "", args)
	r.NoError(err)
	r.Equal(exp, c.args)
}
