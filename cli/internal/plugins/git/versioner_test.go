package git

import (
	"context"
	"testing"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/stretchr/testify/require"
)

func Test_Versioner_BuildVersion(t *testing.T) {
	r := require.New(t)

	vr := &commandRunner{
		stdout: "123",
	}

	v := &Versioner{
		pluginsFn: plugins.Plugins{
			vr,
		}.ScopedPlugins,
	}

	s, err := v.BuildVersion(context.Background(), "")
	r.NoError(err)
	r.Equal(vr.stdout, s)
	r.NotNil(vr.cmd)
	r.Equal([]string{"git", "rev-parse", "--short", "HEAD"}, vr.cmd.Args)
}
