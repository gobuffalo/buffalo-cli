package refresh

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/markbates/refresh/refresh"
	"github.com/stretchr/testify/require"
)

func Test_Fixer(t *testing.T) {
	r := require.New(t)

	root, err := ioutil.TempDir("", "")
	r.NoError(err)

	f, err := os.Create(filepath.Join(root, "go.mod"))
	r.NoError(err)
	f.WriteString("module coke")
	r.NoError(f.Close())

	c := &refresh.Configuration{
		BuildTargetPath: ".",
		BuildPath:       "tmp",
		BuildDelay:      400,
		BinaryName:      "",
		EnableColors:    true,
		LogName:         "buffalo",
	}

	configPath := filepath.Join(root, ".buffalo.dev.yml")
	err = c.Dump(configPath)
	r.NoError(err)

	r.Equal(c.BuildTargetPath, ".")

	fx := &Fixer{}
	r.NoError(fx.Fix(context.Background(), root, []string{}))

	r.NoError(c.Load(configPath))

	expectedBuildPath := "." + string(filepath.Separator) + filepath.Join("cmd", "coke")
	r.Equal(c.BuildTargetPath, expectedBuildPath)
	r.Equal(c.BuildDelay, time.Duration(400))
}

func Test_Fixer_NoFile(t *testing.T) {
	r := require.New(t)

	root, err := ioutil.TempDir("", "")
	r.NoError(err)

	f, err := os.Create(filepath.Join(root, "go.mod"))
	r.NoError(err)
	f.WriteString("module coke")
	r.NoError(f.Close())

	fx := &Fixer{}
	r.NoError(fx.Fix(context.Background(), root, []string{}))

	configPath := filepath.Join(root, ".buffalo.dev.yml")
	c := &refresh.Configuration{}
	r.NoError(c.Load(configPath))

	relative, err := filepath.Rel(root, filepath.Join(root, "cmd", "coke"))

	r.NoError(err)
	r.Equal(c.BuildTargetPath, "."+string(filepath.Separator)+relative)
	r.Equal(c.BuildDelay, time.Duration(200))
}
