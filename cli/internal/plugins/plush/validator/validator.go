package validator

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/build"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plush"
)

var _ build.BeforeBuilder = &Validator{}
var _ plugins.Plugin = Validator{}

type Validator struct{}

func (b Validator) PluginName() string {
	return "plush/validator"
}

func (b *Validator) BeforeBuild(ctx context.Context, root string, args []string) error {
	root = filepath.Join(root, "templates")
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		base := filepath.Base(path)
		if !strings.Contains(base, ".plush") {
			return nil
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		if _, err = plush.Parse(string(b)); err != nil {
			return fmt.Errorf("could not parse %s: %v", path, err)
		}
		return nil
	})
	return plugins.Wrap(b, err)
}
