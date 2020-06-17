package docker

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Generator(t *testing.T) {
	r := require.New(t)
	generator := &Generator{}

	wd, err := os.Getwd()
	r.NoError(err)

	defer os.Chdir(wd)

	root := os.TempDir()
	os.Chdir(root)

	err = generator.Newapp(context.Background(), root, "app name", []string{})
	r.NoError(err)

	b, err := ioutil.ReadFile(filepath.Join(root, "Dockerfile"))
	r.NoError(err)

	r.Contains(string(b), "multi-stage")
	r.Contains(string(b), "FROM alpine")

	generator.style = "standard"
	err = generator.Newapp(context.Background(), root, "app name", []string{})

	b, err = ioutil.ReadFile(filepath.Join(root, "Dockerfile"))
	r.NoError(err)

	r.NotContains(string(b), "multi-stage")
	r.NotContains(string(b), "FROM alpine")
}

func Test_templateFile(t *testing.T) {
	r := require.New(t)

	generator := &Generator{}

	tcases := []struct {
		style    string
		filename string
	}{
		{"", "Dockerfile.multistage"},
		{"standard", "Dockerfile.standard"},
		{"multistage", "Dockerfile.multistage"},
	}

	for index, tcase := range tcases {
		generator.style = tcase.style
		f, err := generator.templateFile()

		r.NoError(err)
		r.Contains(f.Path().String(), tcase.filename, "filename on %v", index)
	}
}

func Test_hasWebpack(t *testing.T) {
	r := require.New(t)
	root := os.TempDir()
	defer os.Remove(filepath.Join(root, "webpack.config.js"))

	generator := &Generator{}
	r.False(generator.hasWebpack(root))

	err := ioutil.WriteFile(filepath.Join(root, "webpack.config.js"), []byte("{}"), 0644)
	r.NoError(err)
	r.True(generator.hasWebpack(root))
}

func Test_tool(t *testing.T) {
	r := require.New(t)
	root := os.TempDir()
	defer os.Remove(filepath.Join(root, "yarn.lock"))

	generator := &Generator{}
	r.Equal(generator.tool(root), "npm")

	err := ioutil.WriteFile(filepath.Join(root, "yarn.lock"), []byte("{}"), 0644)
	r.NoError(err)
	r.Equal(generator.tool(root), "yarn")
}
