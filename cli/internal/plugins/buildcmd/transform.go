package buildcmd

import (
	"context"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/here"
	"github.com/gobuffalo/here/there"
	"golang.org/x/tools/go/ast/astutil"
)

const mainBuildFile = "main.build.go"

type mainFile struct {
	plugins.IO
	plugins func() plugins.Plugins
}

func (mainFile) Name() string {
	return "main"
}

func (bc *mainFile) Version(ctx context.Context, root string) (string, error) {
	versions := map[string]string{
		"time": time.Now().Format(time.RFC3339),
	}
	m := func() (string, error) {
		b, err := json.Marshal(versions)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}

	if bc.plugins == nil {
		return m()
	}

	for _, p := range bc.plugins() {
		bv, ok := p.(BuildVersioner)
		if !ok {
			continue
		}
		plugins.SetIO(bc, p)
		s, err := bv.BuildVersion(ctx, root)
		if err != nil {
			return "", err
		}
		if len(s) == 0 {
			continue
		}
		versions[p.Name()] = strings.TrimSpace(s)
	}
	return m()
}

func (bc *mainFile) generateNewMain(info here.Info, version string) error {
	fmt.Println("version --> ", version)
	bt := struct {
		BuildTime    string
		BuildVersion string
	}{
		BuildVersion: strconv.Quote(version),
		BuildTime:    strconv.Quote(time.Now().Format(time.RFC3339)),
	}

	t, err := template.New(mainBuildFile).Parse(mainBuildTmpl)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(info.Dir, mainBuildFile))
	if err != nil {
		return err
	}
	defer f.Close()

	if err := t.Execute(f, bt); err != nil {
		return err
	}
	return nil
}

func (bc *mainFile) BeforeBuild(ctx context.Context, args []string) error {
	info, err := there.Current()
	if err != nil {
		return err
	}

	err = bc.renameMain(info, "main", "originalMain")
	if err != nil {
		return err
	}

	version, err := bc.Version(ctx, info.Dir)
	if err != nil {
		return err
	}

	if err := bc.generateNewMain(info, version); err != nil {
		return err
	}
	return nil
}

var _ AfterBuilder = &mainFile{}

func (bc *mainFile) AfterBuild(ctx context.Context, args []string, err error) error {
	info, err := there.Current()
	if err != nil {
		return err
	}
	os.RemoveAll(filepath.Join(info.Dir, mainBuildFile))
	err = bc.renameMain(info, "originalMain", "main")
	if err != nil {
		return err
	}

	return nil
}

func (bc *mainFile) renameMain(info here.Info, from string, to string) error {
	if info.Name != "main" {
		return fmt.Errorf("module %s is not a main", info.Name)
	}
	fmt.Println(from, "-->", to)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, info.Dir, nil, 0)
	if err != nil {
		return err
	}

	for _, p := range pkgs {
		for x, f := range p.Files {

			err := func() error {
				var reprint bool
				pre := func(c *astutil.Cursor) bool {
					n := c.Name()
					if n != "Decls" {
						return true
					}

					fd, ok := c.Node().(*ast.FuncDecl)
					if !ok {
						return true
					}

					n = fd.Name.Name
					if n != from {
						return true
					}

					fd.Name = ast.NewIdent(to)
					c.Replace(fd)
					reprint = true
					return true
				}

				res := astutil.Apply(f, pre, nil)
				if !reprint {
					return nil
				}

				f, err := os.Create(x)
				if err != nil {
					return err
				}
				defer f.Close()
				err = printer.Fprint(f, fset, res)
				if err != nil {
					return err
				}
				return nil
			}()

			if err != nil {
				return err
			}
		}
	}
	return nil
}
