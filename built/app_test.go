package built

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo/runtime"
	"github.com/stretchr/testify/require"
)

type plugger []plugins.Plugin

func (p plugger) Plugins() []plugins.Plugin {
	return []plugins.Plugin(p)
}

func Test_App(t *testing.T) {
	r := require.New(t)

	ctx := context.Background()

	stdout := &bytes.Buffer{}

	var res bool
	app := &App{
		IO:      plugins.NewIO(),
		Plugger: plugger{},
		OriginalMain: func() {
			res = true
		},
		BuildVersion: "xxx",
	}

	app.SetStdout(stdout)

	var args []string
	err := app.Main(ctx, args)
	r.NoError(err)
	r.True(res)

	r.Equal(runtime.Build().Version, "xxx")

	args = []string{"version"}
	err = app.Main(ctx, args)
	r.NoError(err)

	s := strings.TrimSpace(stdout.String())
	r.Equal(s, "xxx")

	var fl bool
	app.Fallthrough = func(ctx context.Context, args []string) error {
		fl = true
		return nil
	}

	args = []string{"foo"}
	err = app.Main(ctx, args)
	r.NoError(err)
	r.True(fl)
}
