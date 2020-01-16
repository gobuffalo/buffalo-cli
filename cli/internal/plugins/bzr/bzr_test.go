package bzr

import (
	"bytes"
	"context"
	"testing"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/stretchr/testify/require"
)

// testVersionRunner is a custom version runner used for testing purposes.
type testVersionRunner struct {
	resultError   error
	resultVersion string
}

func (tv *testVersionRunner) ToolAvailable() (bool, error) {
	return true, nil
}

func (tv *testVersionRunner) RunVersionCommand(ctx context.Context, bb *bytes.Buffer) error {
	bb.Write([]byte(tv.resultVersion))
	return tv.resultError
}

func Test_Bzr_Generalities(t *testing.T) {
	r := require.New(t)
	b := Versioner{}

	r.Equal("bzr", b.Name(), "Name should be bzr")
	r.Equal("Provides bzr related hooks to Buffalo applications.", b.Description(), "Description does not match")
}

func Test_Bzr_BuildVersion(t *testing.T) {

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

	r.Equal([]string{"bzr", "revno"}, vr.cmd.Args)
}
