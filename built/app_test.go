package built

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/stretchr/testify/require"
)

func Test_App_No_Args(t *testing.T) {
	r := require.New(t)

	var res bool
	app := &App{
		OriginalMain: func() {
			res = true
		},
	}

	var args []string
	ctx := context.Background()
	err := app.Main(ctx, args)
	r.NoError(err)
	r.True(res)
}

func Test_App_No_Args_Fallthrough(t *testing.T) {
	r := require.New(t)

	var res bool
	app := &App{
		Fallthrough: func(ctx context.Context, args []string) error {
			res = true
			return nil
		},
	}

	var args []string
	ctx := context.Background()
	err := app.Main(ctx, args)
	r.NoError(err)
	r.True(res)
}

func Test_App_With_Args_Fallthrough(t *testing.T) {
	r := require.New(t)

	var res bool
	app := &App{
		Fallthrough: func(ctx context.Context, args []string) error {
			res = true
			return nil
		},
	}

	ctx := context.Background()
	err := app.Main(ctx, []string{"lee", "majors"})
	r.NoError(err)
	r.True(res)
}

func Test_App_Init_Plugins(t *testing.T) {
	r := require.New(t)

	var res bool
	var pres bool

	fn := func(ctx context.Context, args []string) error {
		pres = true
		return nil
	}

	plugs := plugins.Plugins{
		WithIniter(plugins.Background(""), fn),
	}

	app := &App{
		OriginalMain: func() {
			res = true
		},
		Plugger: plugs,
	}

	var args []string
	ctx := context.Background()
	err := app.Main(ctx, args)
	r.NoError(err)
	r.True(res)
	r.True(pres)
}

func Test_App_Init_Plugins_Error(t *testing.T) {
	r := require.New(t)

	var res bool
	var pres bool
	exp := fmt.Errorf("boom")
	fn := func(ctx context.Context, args []string) error {
		return exp
	}

	plugs := plugins.Plugins{
		WithIniter(plugins.Background(""), fn),
	}

	app := &App{
		OriginalMain: func() {
			res = true
		},
		Plugger: plugs,
	}

	var args []string
	ctx := context.Background()
	err := app.Main(ctx, args)
	r.Error(err)
	r.Equal(exp, err)
	r.False(res)
	r.False(pres)
}

func Test_App_Version(t *testing.T) {
	r := require.New(t)

	app := &App{
		IO:           plugins.NewIO(),
		BuildVersion: "v1",
	}
	stdout := &bytes.Buffer{}
	app.SetStdout(stdout)

	ctx := context.Background()
	err := app.Main(ctx, []string{"version"})
	r.NoError(err)

	s := strings.TrimSpace(stdout.String())
	r.Equal("v1", s)
}
