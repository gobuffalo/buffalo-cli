package clifix

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/fix"
	"github.com/gobuffalo/here"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
)

var _ plugins.Plugin = &Fixer{}
var _ plugcmd.Namer = &Fixer{}
var _ fix.Fixer = &Fixer{}

//Fixer creates the cli file at cmd/buffalo/main.go if it doesn't exist.
type Fixer struct {
}

//PluginName for this cli fixer
func (*Fixer) PluginName() string {
	return "cli/fixer"
}

//CmdName for this cli fixer
func (*Fixer) CmdName() string {
	return "cli"
}

//Fix will be invoked when buffalo fix is called, it creates cmd/buffalo/main.go
//with tmplMain if it doesn't exist.
func (fixer *Fixer) Fix(ctx context.Context, root string, args []string) error {
	info, err := here.Dir(root)
	if err != nil {
		return err
	}

	x := filepath.Join(root, "cmd", "buffalo")
	mm := map[string]string{
		filepath.Join(x, "main.go"): tmplMain,
	}

	_, err = os.Stat(filepath.Join(x, "main.go"))
	if err == nil {
		fmt.Println("cmd/buffalo/main.go already exist,s no need to fix it")
		return err
	}

	for fp, body := range mm {
		if err := os.MkdirAll(filepath.Dir(fp), 0755); err != nil {
			return err
		}

		if _, err := os.Stat(fp); err == nil {
			return fmt.Errorf("%s already exists", fp)
		}

		f, err := os.Create(fp)
		if err != nil {
			return err
		}

		body = strings.TrimSpace(body)
		tmpl, err := template.New(fp).Parse(body)
		if err != nil {
			return err
		}

		err = tmpl.Execute(f, struct {
			Name       string
			ImportPath string
		}{
			ImportPath: info.Module.Path,
			Name:       path.Base(info.Module.Path),
		})

		if err != nil {
			return err
		}

		if err := f.Close(); err != nil {
			return err
		}
	}

	return nil
}

const tmplMain = `
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gobuffalo/buffalo-cli/v2/cli"
)

func main() {
	fmt.Print("~~~~ Using coke/cmd/buffalo ~~~\n\n")

	ctx := context.Background()
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	buffalo, err := cli.New()
	if err != nil {
		log.Fatal(err)
	}

	// append your plugins here
	// buffalo.Plugins = append(buffalo.Plugins, ...)

	err = buffalo.Main(ctx, pwd, os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}
`
