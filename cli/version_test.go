package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	bufcli "github.com/gobuffalo/buffalo-cli"
	"github.com/stretchr/testify/require"
)

func Test_Buffalo_Version(t *testing.T) {
	r := require.New(t)

	buffalo, err := New()
	r.NoError(err)

	bb := &bytes.Buffer{}
	buffalo.Stdout = bb

	err = buffalo.Main(context.Background(), []string{"version"})
	r.NoError(err)

	out := strings.TrimSpace(bb.String())
	r.Equal(out, bufcli.Version)
}

func Test_Buffalo_Version_JSON(t *testing.T) {
	r := require.New(t)

	buffalo, err := New()
	r.NoError(err)

	bb := &bytes.Buffer{}
	buffalo.Stdout = bb

	err = buffalo.Main(context.Background(), []string{"version", "--json"})
	r.NoError(err)

	out := strings.TrimSpace(bb.String())
	r.Contains(out, "version")
	r.Contains(out, bufcli.Version)
}
