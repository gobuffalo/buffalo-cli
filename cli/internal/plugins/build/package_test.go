package build

import (
	"context"
	"testing"

	"github.com/gobuffalo/plugins"
	"github.com/stretchr/testify/require"
)

func Test_Cmd_Package(t *testing.T) {
	r := require.New(t)

	pkg := &packager{
		files: []string{"A"},
	}
	pf := &packFiler{
		files: []string{"B"},
	}

	plugs := plugins.Plugins{
		pkg,
		pf,
		&bladeRunner{},
	}

	bc := &Cmd{}
	bc.WithPlugins(func() []plugins.Plugin {
		return plugs
	})

	err := bc.Main(context.Background(), ".", nil)
	r.NoError(err)

	r.Len(pkg.files, 2)
	r.Equal([]string{"A", "B"}, pkg.files)
}
