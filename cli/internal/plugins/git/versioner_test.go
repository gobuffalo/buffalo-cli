package git

import (
	"context"
	"testing"

	"github.com/gobuffalo/plugins"
	"github.com/stretchr/testify/require"
)

func Test_Versioner_BuildVersion(t *testing.T) {
	r := require.New(t)

	vr := &commandRunner{
		stdout: "123",
	}

	v := &Versioner{
		pluginsFn: func() []plugins.Plugin {
			return []plugins.Plugin{
				vr,
			}
		},
	}

	s, err := v.BuildVersion(context.Background(), "")
	r.NoError(err)
	r.Equal(vr.stdout, s)
	r.NotNil(vr.cmd)
	r.Equal([]string{"git", "rev-parse", "--short", "HEAD"}, vr.cmd.Args)
}
