package buildcmd

import (
	"bytes"
	"context"
	"os"
	"path"
	"testing"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/stretchr/testify/require"
)

func Test_mainFile_Version(t *testing.T) {
	r := require.New(t)

	ref := newRef(t, "ref")
	defer ref.Close()
	os.Chdir(ref.Root)

	bc := &mainFile{}

	ctx := context.Background()

	s, err := bc.Version(ctx, ref.Root)
	r.NoError(err)
	r.Contains(s, `"time":`)

	bc.PluginsFn = func() []plugins.Plugin {
		return plugins.Plugins{
			&buildVersioner{version: "v1"},
		}
	}

	s, err = bc.Version(ctx, ref.Root)
	r.NoError(err)
	r.Contains(s, `"time":`)
	r.Contains(s, `"buildVersioner":"v1"`)
}

func Test_mainFile_generateNewMain(t *testing.T) {
	r := require.New(t)

	ref := newRef(t, "ref")
	defer ref.Close()
	os.Chdir(ref.Root)

	plugs := plugins.Plugins{
		&buildImporter{
			imports: []string{
				path.Join(ref.ImportPath, "actions"),
			},
		},
	}
	bc := &mainFile{
		PluginsFn: plugs.WithPlugins,
	}

	ctx := context.Background()
	bb := &bytes.Buffer{}
	err := bc.generateNewMain(ctx, ref.Info, "v1", bb)
	r.NoError(err)

	out := bb.String()
	r.Contains(out, `appcli "github.com/markbates/coke/cli"`)
	r.Contains(out, `_ "github.com/markbates/coke/actions"`)
	r.Contains(out, `appcli.Buffalo`)
	r.Contains(out, `originalMain`)
}

func Test_mainFile_generateNewMain_noCli(t *testing.T) {
	r := require.New(t)

	ref := newRef(t, "nocli")
	defer ref.Close()
	os.Chdir(ref.Root)

	plugs := plugins.Plugins{
		&buildImporter{
			imports: []string{
				path.Join(ref.ImportPath, "actions"),
			},
		},
	}
	bc := &mainFile{
		PluginsFn: plugs.WithPlugins,
	}

	ctx := context.Background()
	bb := &bytes.Buffer{}
	err := bc.generateNewMain(ctx, ref.Info, "v1", bb)
	r.NoError(err)

	out := bb.String()
	r.NotContains(out, `appcli "github.com/markbates/coke/cli"`)
	r.NotContains(out, `appcli.Buffalo`)
	r.Contains(out, `_ "github.com/markbates/coke/actions"`)
	r.Contains(out, `originalMain`)
	r.Contains(out, `cb.Main`)
}
