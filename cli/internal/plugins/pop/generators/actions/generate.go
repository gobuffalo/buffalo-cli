package actions

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"

	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/here"
	"github.com/gobuffalo/meta/v2"
	"github.com/gobuffalo/plugins"
)

func (mg *Generator) GenerateResourceActions(ctx context.Context, root string, args []string) error {
	if len(args) == 0 {
		err := fmt.Errorf("you must specify a resource")
		return plugins.Wrap(mg, err)
	}

	info, err := here.Dir(root)
	if err != nil {
		return plugins.Wrap(mg, err)
	}

	resourceName := args[0]

	modelName := mg.modelName
	if len(modelName) == 0 {
		modelName = resourceName
	}

	modelsPkg := mg.modelsPkg
	if len(modelsPkg) == 0 {
		modelsPkg = path.Join(info.ImportPath, "models")
	}

	modelsPkgSel := mg.modelsPkgSel
	importName := modelsPkgSel
	if len(modelsPkgSel) == 0 {
		modelsPkgSel = path.Base(modelsPkg)
		importName = ""
	}

	app, err := meta.New(info)
	if err != nil {
		return plugins.Wrap(mg, err)
	}

	pres := struct {
		AsWeb        bool
		ImportName   string
		Model        name.Ident
		ModelsPkg    string
		ModelsPkgSel string
		Name         name.Ident
	}{
		AsWeb:        app.As["web"],
		ImportName:   importName,
		Model:        name.New(modelName),
		ModelsPkg:    modelsPkg,
		ModelsPkgSel: modelsPkgSel,
		Name:         name.New(resourceName),
	}

	t, err := template.New(resourceName).Parse(actionsTmpl)
	if err != nil {
		return plugins.Wrap(mg, err)
	}

	p := fmt.Sprintf("%s.go", pres.Name.Folder().Pluralize())
	fp := filepath.Join(root, "actions", p)

	if err := os.MkdirAll(filepath.Dir(fp), 0755); err != nil {
		return plugins.Wrap(mg, err)
	}

	f, err := os.Create(fp)
	if err != nil {
		return plugins.Wrap(mg, err)
	}
	defer f.Close()

	if err := t.Execute(f, pres); err != nil {
		return plugins.Wrap(mg, err)
	}

	return nil
}
