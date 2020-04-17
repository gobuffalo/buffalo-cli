package build

import (
	"context"
	"testing"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/build/buildtest"
	"github.com/gobuffalo/plugins"
	"github.com/stretchr/testify/require"
)

func Test_Cmd_Package(t *testing.T) {
	r := require.New(t)

	exp := []string{"foo.go", "bar.go"}
	pf := func(ctx context.Context, root string) ([]string, error) {
		return exp, nil
	}

	var act []string
	pr := func(ctx context.Context, root string, files []string) error {
		act = files
		return nil
	}

	plugs := plugins.Plugins{
		buildtest.Packager(pr),
		buildtest.PackFiler(pf),
	}

	bc := &Cmd{}
	bc.WithPlugins(func() []plugins.Plugin {
		return plugs
	})

	err := bc.Main(context.Background(), ".", nil)
	r.NoError(err)

	r.Equal(exp, act)
}
