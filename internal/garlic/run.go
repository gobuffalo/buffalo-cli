package garlic

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/v2/cli"
	"github.com/gobuffalo/here"
	"github.com/gobuffalo/plugins/plugio"
)

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

	bargs := []string{"run", "-v", "./cmd/buffalo"}
	bargs = append(bargs, args...)

	cmd := exec.CommandContext(ctx, "go", bargs...)
	cmd.Stdin = plugio.Stdin()
	cmd.Stdout = plugio.Stdout()
	cmd.Stderr = plugio.Stderr()
	return cmd.Run()
}
