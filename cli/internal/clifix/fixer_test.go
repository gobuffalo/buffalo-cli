package clifix

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Fixer_Fix(t *testing.T) {
	r := require.New(t)

	dir, err := ioutil.TempDir("", "")
	r.NoError(err)

	f, err := os.Create(filepath.Join(dir, "go.mod"))
	r.NoError(err)
	f.WriteString("module coke")
	r.NoError(f.Close())

	ctx := context.Background()
	var args []string

	fixer := &Fixer{}
	err = fixer.Fix(ctx, dir, args)
	r.NoError(err)

	root := filepath.Join(dir, "cmd", "buffalo")

	_, err = os.Stat(root)
	r.NoError(err)

	fp := filepath.Join(root, "main.go")
	b, err := ioutil.ReadFile(fp)
	r.NoError(err)
	r.Contains(string(b), `coke/cmd/buffalo`)

}

func Test_Fixer_FileExists(t *testing.T) {
	r := require.New(t)

	dir, err := ioutil.TempDir("", "")
	r.NoError(err)

	f, err := os.Create(filepath.Join(dir, "go.mod"))
	r.NoError(err)
	f.WriteString("module pagano")
	r.NoError(f.Close())

	ctx := context.Background()
	var args []string

	cliFolder := filepath.Join(dir, "cmd", "buffalo")
	err = os.MkdirAll(cliFolder, 0755)
	r.NoError(err)

	fixer := &Fixer{}
	err = fixer.Fix(ctx, dir, args)
	r.NoError(err)
}
