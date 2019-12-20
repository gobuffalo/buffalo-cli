package buildcmd

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Main(t *testing.T) {
	r := require.New(t)

	ref := newRef(t, "ref")
	defer ref.Close()
	os.Chdir(ref.Root)

	bc := &BuildCmd{}

	ctx := context.Background()
	var args []string

	err := bc.Main(ctx, args)
	r.NoError(err)
}
