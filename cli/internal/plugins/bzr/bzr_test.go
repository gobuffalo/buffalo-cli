package bzr

import (
	"context"
	"os/exec"
	"testing"

	"github.com/gobuffalo/plugins"
	"github.com/stretchr/testify/require"
)

func Test_Bzr_Generalities(t *testing.T) {
	r := require.New(t)
	b := Versioner{}

	r.Equal("bzr", b.PluginName(), "Name should be bzr")
	r.Equal("Provides bzr related hooks to Buffalo applications.", b.Description(), "Description does not match")
}

func Test_Bzr_BuildVersion(t *testing.T) {

	r := require.New(t)

	var act []string
	fn := func(ctx context.Context, root string, cmd *exec.Cmd) error {
		act = cmd.Args
		if cmd.Stdout == nil {
			r.FailNow("expected stdout not to be nil")
		}
		cmd.Stdout.Write([]byte("42"))
		return nil
	}

	v := &Versioner{
		pluginsFn: func() []plugins.Plugin {
			return []plugins.Plugin{
				runner(fn),
			}
		},
	}

	s, err := v.BuildVersion(context.Background(), "")
	r.NoError(err)
	r.Equal("42", s)
	r.Equal([]string{"bzr", "revno"}, act)
}
