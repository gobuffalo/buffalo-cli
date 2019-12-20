package buildcmd

import (
	"context"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/here"
	"golang.org/x/tools/go/ast/astutil"
)

const mainBuildFile = "main.build.go"

var _ AfterBuilder = &mainFile{}
var _ BeforeBuilder = &mainFile{}
var _ plugins.Plugin = &mainFile{}
var _ plugins.PluginNeeder = &mainFile{}
var _ plugins.PluginScoper = &mainFile{}
var _ plugprint.Hider = &mainFile{}

type mainFile struct {
	pluginsFn plugins.PluginFeeder
}

func (bc *mainFile) WithPlugins(f plugins.PluginFeeder) {
	bc.pluginsFn = f
}

func (mainFile) HidePlugin() {}

func (mainFile) Name() string {
	return "main"
}

func (bc *mainFile) ScopedPlugins() []plugins.Plugin {
	if bc.pluginsFn == nil {
		return nil
	}
	return bc.pluginsFn()
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

	for _, p := range bc.ScopedPlugins() {
		bv, ok := p.(BuildVersioner)
		if !ok {
			continue
		}
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

func (bc *mainFile) generateNewMain(ctx context.Context, info here.Info, version string, ws ...io.Writer) error {
	fmt.Println("version --> ", version)

	var imports []string
	for _, p := range bc.ScopedPlugins() {
		bi, ok := p.(BuildImporter)
		if !ok {
			continue
		}
		i, err := bi.BuildImports(ctx, info.Root)
		if err != nil {
			return err
		}
		imports = append(imports, i...)
	}

	if i, err := here.Dir(filepath.Join(info.Dir, "actions")); err == nil {
		imports = append(imports, i.ImportPath)
	}

	sort.Strings(imports)

	bt := struct {
		BuildTime       string
		BuildVersion    string
		Imports         []string
		Info            here.Info
		WithFallthrough bool
	}{
		BuildTime:    strconv.Quote(time.Now().Format(time.RFC3339)),
		BuildVersion: strconv.Quote(version),
		Imports:      imports,
		Info:         info,
	}

	bt.WithFallthrough = func() bool {
		c := exec.CommandContext(ctx, "go", "doc", path.Join(info.ImportPath, "cli")+".Buffalo")
		err := c.Run()
		return err == nil
	}()

	t, err := template.New(mainBuildFile).Parse(mainBuildTmpl)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(info.Dir, mainBuildFile))
	if err != nil {
		return err
	}
	defer f.Close()

	if err := t.Execute(io.MultiWriter(append(ws, f)...), bt); err != nil {
		return err
	}
	return nil
}

func (bc *mainFile) BeforeBuild(ctx context.Context, args []string) error {
	info, err := here.Current()
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

	if err := bc.generateNewMain(ctx, info, version, os.Stdout); err != nil {
		return err
	}
	return nil
}

var _ AfterBuilder = &mainFile{}

func (bc *mainFile) AfterBuild(ctx context.Context, args []string, err error) error {
	info, err := here.Current()
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
