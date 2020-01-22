package i18n

import (
	"context"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/gobuffalo/flect/name"
)

var _ plugins.PluginNeeder = &flasher{}
var _ plugins.PluginScoper = &flasher{}

type flasher struct {
	pluginsFn plugins.PluginFeeder
	model     name.Ident
	fset      *token.FileSet
}

func (flash *flasher) WithPlugins(f plugins.PluginFeeder) {
	flash.pluginsFn = f
}

func (flash *flasher) ScopedPlugins() []plugins.Plugin {
	if flash.pluginsFn == nil {
		return []plugins.Plugin{}
	}
	return flash.pluginsFn()
}

func (flash *flasher) namedWriter(ctx context.Context, filename string) (io.Writer, error) {
	for _, p := range flash.ScopedPlugins() {
		if fw, ok := p.(NamedWriter); ok {
			return fw.NamedWriter(ctx, filename)
		}
	}
	return os.Create(filename)
}

func (flash *flasher) Flash(ctx context.Context, root string, model name.Ident) error {
	flash.model = model
	flash.fset = token.NewFileSet()
	pkgs, err := parser.ParseDir(flash.fset, filepath.Join(root, "actions"), nil, 0)
	if err != nil {
		return err
	}

	pkg, ok := pkgs["actions"]
	if !ok {
		return nil
	}

	for n, f := range pkg.Files {
		if err := flash.file(ctx, n, f); err != nil {
			return err
		}
	}
	return nil
}

func (flash *flasher) file(ctx context.Context, filename string, f *ast.File) error {
	var flashed bool
	for _, d := range f.Decls {
		fnd, ok := d.(*ast.FuncDecl)
		if !ok {
			continue
		}
		recv := fnd.Recv
		if recv == nil {
			continue
		}
		for _, rv := range recv.List {
			id, ok := rv.Type.(*ast.Ident)
			if !ok {
				continue
			}

			if id.Name != fmt.Sprintf("%sResource", flash.model.Resource()) {
				continue
			}

			if err := flash.block(f, fnd); err != nil {
				return err
			}
			flashed = true
			continue
		}
	}

	if flashed {
		w, err := flash.namedWriter(ctx, filename)
		if err != nil {
			return err
		}
		if c, ok := w.(io.Closer); ok {
			defer c.Close()
		}
		return format.Node(w, flash.fset, f)
	}

	return nil
}

func (flash *flasher) block(f *ast.File, fnd *ast.FuncDecl) error {
	body := fnd.Body
	if body == nil {
		return nil
	}
	const msg = `T.Translate(c, "%s.%s.success")`

	var state string
	switch fnd.Name.Name {
	case "Create":
		state = "created"
	case "Update":
		state = "updated"
	case "Destroy":
		state = "destroyed"
	default:
		return nil
	}

	for _, s := range body.List {
		exs, ok := s.(*ast.ExprStmt)
		if !ok {
			continue
		}
		ce, ok := exs.X.(*ast.CallExpr)
		if !ok {
			continue
		}
		pse, ok := ce.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		if pse.Sel.Name != "Add" || len(ce.Args) != 2 {
			continue
		}
		arg := ce.Args[1]
		bl, ok := arg.(*ast.BasicLit)
		if !ok {
			continue
		}
		bl.Value = fmt.Sprintf(msg, flash.model.VarCaseSingle(), state)
		ce.Args[1] = bl
	}
	return nil
}
