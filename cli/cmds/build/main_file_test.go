package build

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/build/buildtest"
	"github.com/gobuffalo/plugins"
	"github.com/stretchr/testify/require"
)

func Test_MainFile_Version(t *testing.T) {
	r := require.New(t)

	bc := &MainFile{}

	ctx := context.Background()

	s, err := bc.Version(ctx, "")
	r.NoError(err)
	r.Contains(s, `"time":`)

	fn := func(ctx context.Context, root string) (string, error) {
		return "v1", nil
	}

	p := buildtest.Versioner(fn)
	bc.pluginsFn = func() []plugins.Plugin {
		return plugins.Plugins{p}
	}

	s, err = bc.Version(ctx, "")
	r.NoError(err)
	r.Contains(s, `"time":`)
	r.Contains(s, fmt.Sprintf(`"%s":"v1"`, p.PluginName()))
}

func Test_MainFile_generateNewMain(t *testing.T) {
	r := require.New(t)

	ref := newRef(t, "")
	defer os.RemoveAll(filepath.Join(ref.Dir, mainBuildFile))

	fn := func(ctx context.Context, root string) ([]string, error) {
		return []string{
			path.Join(ref.ImportPath, "actions"),
		}, nil
	}
	plugs := plugins.Plugins{
		buildtest.Importer(fn),
	}
	bc := &MainFile{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
		withFallthroughFn: func() bool { return true },
	}

	ctx := context.Background()
	bb := &bytes.Buffer{}
	err := bc.generateNewMain(ctx, ref, "v1", bb)
	r.NoError(err)

	out := bb.String()
	r.Contains(out, `appcli "github.com/markbates/coke/cli"`)
	r.Contains(out, `_ "github.com/markbates/coke/actions"`)
	r.Contains(out, `appcli.Buffalo`)
	r.Contains(out, `originalMain`)
}

func Test_MainFile_generateNewMain_noCli(t *testing.T) {
	r := require.New(t)

	ref := newRef(t, "")
	defer os.RemoveAll(filepath.Join(ref.Dir, mainBuildFile))

	fn := func(ctx context.Context, root string) ([]string, error) {
		return []string{
			path.Join(ref.ImportPath, "actions"),
		}, nil
	}
	plugs := plugins.Plugins{
		buildtest.Importer(fn),
	}
	bc := &MainFile{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
		withFallthroughFn: func() bool { return false },
	}

	bb := &bytes.Buffer{}
	err := bc.generateNewMain(context.Background(), ref, "v1", bb)
	r.NoError(err)

	out := bb.String()
	r.NotContains(out, `appcli "github.com/markbates/coke/cli"`)
	r.NotContains(out, `appcli.Buffalo`)
	r.Contains(out, `_ "github.com/markbates/coke/actions"`)
	r.Contains(out, `originalMain`)
	r.Contains(out, `cb.Main`)
}
