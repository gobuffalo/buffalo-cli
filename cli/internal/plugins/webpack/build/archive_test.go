package build

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Builder_archive(t *testing.T) {
	r := require.New(t)

	ctx := context.Background()

	info := tempApp(t, map[string]string{})

	ap := filepath.Join(info.Dir, "public", "webpack")
	os.MkdirAll(ap, 0755)
	err := ioutil.WriteFile(filepath.Join(ap, "app.css"), []byte(""), 0644)
	r.NoError(err)

	defer os.RemoveAll(info.Dir)

	pwd, err := os.Getwd()
	r.NoError(err)
	defer os.Chdir(pwd)

	os.Chdir(info.Dir)

	bc := &Builder{
		extract: true,
	}

	err = bc.archive(ctx, info.Dir, nil)
	r.NoError(err)

	_, err = os.Stat(filepath.Join(info.Dir, "bin", "webpack.zip"))
	r.NoError(err)
}

func Test_Builder_archive_extract_to(t *testing.T) {
	r := require.New(t)

	ctx := context.Background()

	info := tempApp(t, map[string]string{})

	ap := filepath.Join(info.Dir, "public", "webpack")
	os.MkdirAll(ap, 0755)
	err := ioutil.WriteFile(filepath.Join(ap, "app.css"), []byte(""), 0644)
	r.NoError(err)

	defer os.RemoveAll(info.Dir)

	pwd, err := os.Getwd()
	r.NoError(err)
	defer os.Chdir(pwd)

	os.Chdir(info.Dir)

	bc := &Builder{
		extract:   true,
		extractTo: "wolverine",
	}

	err = bc.archive(ctx, info.Dir, nil)
	r.NoError(err)

	_, err = os.Stat(filepath.Join(info.Dir, "wolverine", "webpack.zip"))
	r.NoError(err)
}
