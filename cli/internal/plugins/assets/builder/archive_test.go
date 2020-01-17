package builder

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

	ap := filepath.Join(info.Dir, "public", "assets")
	os.MkdirAll(ap, 0755)
	err := ioutil.WriteFile(filepath.Join(ap, "app.css"), []byte(""), 0644)
	r.NoError(err)

	defer os.RemoveAll(info.Dir)

	pwd, err := os.Getwd()
	r.NoError(err)
	defer os.Chdir(pwd)

	os.Chdir(info.Dir)

	bc := &Builder{
		Extract: true,
	}
	bc.WithHereInfo(info)

	err = bc.archive(ctx, info.Dir, nil)
	r.NoError(err)

	_, err = os.Stat(filepath.Join(info.Dir, "bin", "assets.zip"))
	r.NoError(err)
}

func Test_Builder_archive_extract_to(t *testing.T) {
	r := require.New(t)

	ctx := context.Background()

	info := tempApp(t, map[string]string{})

	ap := filepath.Join(info.Dir, "public", "assets")
	os.MkdirAll(ap, 0755)
	err := ioutil.WriteFile(filepath.Join(ap, "app.css"), []byte(""), 0644)
	r.NoError(err)

	defer os.RemoveAll(info.Dir)

	pwd, err := os.Getwd()
	r.NoError(err)
	defer os.Chdir(pwd)

	os.Chdir(info.Dir)

	bc := &Builder{
		Extract:   true,
		ExtractTo: "wolverine",
	}
	bc.WithHereInfo(info)

	err = bc.archive(ctx, info.Dir, nil)
	r.NoError(err)

	_, err = os.Stat(filepath.Join(info.Dir, "wolverine", "assets.zip"))
	r.NoError(err)
}
