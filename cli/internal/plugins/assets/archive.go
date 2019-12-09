package assets

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/gobuffalo/meta/v2"
	"github.com/markbates/pkger"
)

func (a *Builder) archive(app *meta.App) error {
	if !a.ExtractAssets {
		return nil
	}

	outputDir := filepath.Dir(filepath.Join(app.Info.Root, app.Bin))
	os.MkdirAll(outputDir, 0755)
	target := filepath.Join(outputDir, "assets.zip")
	source := filepath.Join(app.Info.Root, "public", "assets")

	f, err := pkger.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	archive := zip.NewWriter(f)
	defer archive.Close()

	err = pkger.Walk(pkger.Include(source), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		pt, err := app.Info.Parse(path)
		if err != nil {
			return err
		}
		header.Name = pt.Name

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

		file, err := pkger.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
		if err != nil {
			return err
		}
		archive.Close()
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
