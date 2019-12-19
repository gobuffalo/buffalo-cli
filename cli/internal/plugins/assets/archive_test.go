package assets

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/stretchr/testify/require"
)

func Test_Builder_archive(t *testing.T) {
	r := require.New(t)

	info := tempApp(t, map[string]string{
		"build": "echo wolverine",
	})

	ap := filepath.Join(info.Root, "public", "assets")
	os.MkdirAll(ap, 0755)
	err := ioutil.WriteFile(filepath.Join(ap, "app.css"), []byte(""), 0644)
	r.NoError(err)

	defer os.RemoveAll(info.Root)

	pwd, err := os.Getwd()
	r.NoError(err)
	defer os.Chdir(pwd)

	os.Chdir(info.Root)

	bc := &Builder{}
	ctx := context.Background()
	ctx = context.WithValue(ctx, "here.Current", info)

	stdout := &bytes.Buffer{}
	ctx = plugins.WithStdout(ctx, stdout)

	args := []string{"--extract-assets"}

	err = bc.Build(ctx, args)
	r.NoError(err)
	r.Contains(stdout.String(), "wolverine")

	_, err = os.Stat(filepath.Join(info.Root, "bin", "assets.zip"))
	r.NoError(err)
}

func Test_Builder_archive_extract_to(t *testing.T) {
	r := require.New(t)

	info := tempApp(t, map[string]string{
		"build": "echo wolverine",
	})

	ap := filepath.Join(info.Root, "public", "assets")
	os.MkdirAll(ap, 0755)
	err := ioutil.WriteFile(filepath.Join(ap, "app.css"), []byte(""), 0644)
	r.NoError(err)

	defer os.RemoveAll(info.Root)

	pwd, err := os.Getwd()
	r.NoError(err)
	defer os.Chdir(pwd)

	os.Chdir(info.Root)

	bc := &Builder{}
	ctx := context.Background()
	ctx = context.WithValue(ctx, "here.Current", info)

	stdout := &bytes.Buffer{}
	ctx = plugins.WithStdout(ctx, stdout)

	args := []string{"--extract-assets", "--extract-assets-to", "wilco"}

	err = bc.Build(ctx, args)
	r.NoError(err)
	r.Contains(stdout.String(), "wolverine")

	_, err = os.Stat(filepath.Join(info.Root, "wilco", "assets.zip"))
	r.NoError(err)
}
