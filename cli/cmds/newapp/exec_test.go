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
	during := func(ctx context.Context, root string, name string, args []string) error {
		log = append(log, "during", name)
		return nil
	}
	after := func(ctx context.Context, root string, name string, args []string, err error) error {
		log = append(log, "after", name)
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

	err := Execute(plugs, ctx, root, "coke", args)
	r.NoError(err)

	r.Len(log, 6)
	r.Equal([]string{"during", "coke", "during", "coke", "after", "coke"}, log)
}
