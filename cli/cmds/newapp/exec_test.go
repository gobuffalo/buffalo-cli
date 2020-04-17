package newapp

import (
	"context"
	"testing"

	"github.com/gobuffalo/plugins"
	"github.com/stretchr/testify/require"
)

func Test_Execute(t *testing.T) {
	r := require.New(t)

	var log []string
	during := func(ctx context.Context, root string, args []string) error {
		log = append(log, "during")
		return nil
	}
	after := func(ctx context.Context, root string, args []string, err error) error {
		log = append(log, "after")
		return nil
	}

	plugs := []plugins.Plugin{
		newapper(during),
		afternewapper(after),
		newapper(during),
	}

	ctx := context.Background()
	var root string
	var args []string

	err := Execute(plugs, ctx, root, args)
	r.NoError(err)

	r.Len(log, 3)
	r.Equal([]string{"during", "during", "after"}, log)
}
