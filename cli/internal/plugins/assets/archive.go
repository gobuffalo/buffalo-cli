package assets

import (
	"archive/zip"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/here/there"
	"github.com/gobuffalo/meta/v2"
)

func (a *Builder) archive(ctx context.Context, args []string) error {
	if !a.ExtractAssets {
		return nil
	}

	info, err := there.Current()
	if err != nil {
		return err
	}

	app, err := meta.New(info)
	if err != nil {
		return err
	}

	outputDir := filepath.Dir(filepath.Join(app.Info.Root, app.Bin))
	os.MkdirAll(outputDir, 0755)
	target := filepath.Join(outputDir, "assets.zip")
	source := filepath.Join(app.Info.Root, "public", "assets")

	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	archive := zip.NewWriter(f)
	defer archive.Close()
	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(path, source)
		header.Name = strings.TrimPrefix(header.Name, string(filepath.Separator))

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
