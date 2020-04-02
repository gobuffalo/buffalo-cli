package mod

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
)

var _ plugins.Plugin = &Initer{}
var _ plugins.Needer = &Initer{}
var _ plugins.Scoper = &Initer{}

type Initer struct {
	pluginsFn plugins.Feeder
}

func (i *Initer) WithPlugins(f plugins.Feeder) {
	i.pluginsFn = f
}

func (i *Initer) ScopedPlugins() []plugins.Plugin {
	if i.pluginsFn == nil {
		return nil
	}

	var plugs []plugins.Plugin
	for _, p := range i.pluginsFn() {
		switch p.(type) {
		case Stderrer:
			plugs = append(plugs, p)
		case Stdiner:
			plugs = append(plugs, p)
		case Stdouter:
			plugs = append(plugs, p)
		case Requirer:
			plugs = append(plugs, p)
		case Replacer:
			plugs = append(plugs, p)
		}
	}

	return plugs
}

func (i Initer) PluginName() string {
	return "go/mod/init"
}

func (i *Initer) ModInit(ctx context.Context, root string, name string) error {
	plugs := i.ScopedPlugins()

	c := exec.CommandContext(ctx, "go", "mod", "init", name)
	c.Stdout = plugio.Stdout(plugs...)
	c.Stderr = plugio.Stderr(plugs...)
	c.Stdin = plugio.Stdin(plugs...)

	if err := c.Run(); err != nil {
		return err
	}

	c = exec.CommandContext(ctx, "go", "mod", "tidy", "-v")
	c.Stdout = plugio.Stdout(plugs...)
	c.Stderr = plugio.Stderr(plugs...)
	c.Stdin = plugio.Stdin(plugs...)

	if len(plugs) == 0 {
		return nil
	}

	mr := filepath.Join(root, "go.mod")
	b, err := ioutil.ReadFile(mr)
	if err != nil {
		return err
	}

	bb := bytes.NewBuffer(b)

	for _, p := range plugs {
		switch t := p.(type) {
		case Replacer:
			m := t.ModReplace(root)
			fmt.Println(">>>TODO cli/internal/plugins/golang/mod/init.go:84: m ", m)
			for k, v := range m {
				s := fmt.Sprintf("\nreplace %s => %s", k, v)
				bb.WriteString(s)
			}
		}
	}
	fmt.Println(bb.String())

	f, err := os.Open(mr)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write(bb.Bytes())
	return nil
}
