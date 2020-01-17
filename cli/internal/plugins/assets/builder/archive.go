package builder

import (
	"archive/zip"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func (a *Builder) archive(ctx context.Context, root string, args []string) error {
	if !a.Extract {
		return nil
	}

	outputDir := a.ExtractTo
	if len(a.ExtractTo) == 0 {
		outputDir = "bin"
	}
	outputDir = filepath.Join(root, outputDir)
	os.MkdirAll(outputDir, 0755)

	target := filepath.Join(outputDir, "assets.zip")
	source := filepath.Join(root, "public", "assets")

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
