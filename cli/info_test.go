package cli

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

type snow struct {
	io.Writer
}

func (s snow) Name() string {
	return "snow"
}

func (s snow) Info(ctx context.Context, args []string) error {
	s.Write([]byte("informer"))
	return nil
}

func Test_Buffalo_Info(t *testing.T) {
	r := require.New(t)

	buffalo, err := New()
	r.NoError(err)

	bb := &bytes.Buffer{}
	buffalo.Stdout = bb

	buffalo.Plugins = append(buffalo.Plugins, snow{
		Writer: bb,
	})

	err = buffalo.Main(context.Background(), []string{"info"})
	r.NoError(err)

	out := bb.String()
	r.Contains(out, "Buffalo (CLI)")
	r.Contains(out, "informer")
}
