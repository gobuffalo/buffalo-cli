package cli

import (
	"testing"

	"github.com/gobuffalo/plugins"
	"github.com/stretchr/testify/require"
)

func Test_Buffalo_New(t *testing.T) {
	r := require.New(t)

	b, err := New()
	r.NoError(err)
	r.NotNil(b)
	r.NotEmpty(b.Plugins)
}

func Test_Buffalo_SubCommands(t *testing.T) {
	r := require.New(t)

	c := &cp{}
	b := &Buffalo{
		Plugins: plugins.Plugins{
			background("foo"),
			c,
		},
	}
	r.Len(b.Plugins, 2)

	cmds := b.SubCommands()
	r.Len(cmds, 1)
	r.Equal(c, cmds[0])
}
