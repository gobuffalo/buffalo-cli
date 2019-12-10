package golang

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/markbates/pkger"
)

type Templates struct {
}

func (t *Templates) ValidateTemplates(root string) error {
	return pkger.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		base := filepath.Base(path)
		if !strings.Contains(base, ".tmpl") {
			return nil
		}

		f, err := pkger.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		b, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		t := template.New(f.Name())
		if _, err = t.Parse(string(b)); err != nil {
			return fmt.Errorf("could not parse %s: %v", path, err)
		}
		return nil
	})
}

func (t Templates) Name() string {
	return "templates"
}
