package garlic

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/v2/cli"
	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/gobuffalo/here"
	"github.com/markbates/jim"
	"github.com/markbates/safe"
)

type tasker interface {
	Task() *jim.Task
}

func isBuffalo(mod string) bool {

	if _, err := os.Stat(mod); err != nil {
		return false
	}

	b, err := ioutil.ReadFile(mod)
	if err != nil {
		return false
	}

	return bytes.Contains(b, []byte("github.com/gobuffalo/buffalo"))
}

func Run(ctx context.Context, root string, args []string) error {
	info, err := here.Dir(root)
	if err != nil {
		return err
	}

	if !isBuffalo(info.Module.GoMod) {
		// TODO use cli.Buffalo with a set of appropriate
		// plugins for use outside of an app. such as `buffalo new`.
		return fmt.Errorf("%s is not a buffalo app", info.Module)
	}

	main := filepath.Join(root, "cmd", "buffalo")
	if _, err := os.Stat(filepath.Dir(main)); err != nil {
		buff, err := cli.NewFromRoot(root)
		if err != nil {
			return err
		}
		return buff.Main(ctx, root, args)
	}

	bargs := []string{"run", "./cmd/buffalo"}
	bargs = append(bargs, args...)

	cmd := plugins.Cmd(ctx, "go", bargs...)
	err = safe.RunE(func() error {
		// fmt.Println(cmd.Args)
		return cmd.Run()
	})
	if err != nil {
		return err
	}

	return nil

}

func buildTags(ctx context.Context, info here.Info) ([]string, error) {
	var args []string
	dy := filepath.Join(info.Dir, "database.yml")
	if _, err := os.Stat(dy); err != nil {
		return args, nil
	}

	b, err := ioutil.ReadFile(dy)
	if err != nil {
		return nil, err
	}
	if bytes.Contains(b, []byte("sqlite")) {
		args = append(args, "-tags", "sqlite")
	}
	return args, nil
}
