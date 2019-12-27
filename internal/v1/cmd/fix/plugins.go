package fix

import (
	"context"
	"fmt"
	"os"

	"github.com/gobuffalo/buffalo-cli/internal/v1/genny/plugins/install"
	"github.com/gobuffalo/buffalo-cli/internal/v1/plugins"
	"github.com/gobuffalo/buffalo-cli/internal/v1/plugins/plugdeps"
	"github.com/gobuffalo/genny/v2"
	"github.com/markbates/errx"
)

// Plugins will fix plugins between releases
func Plugins(r *Runner) error {
	fmt.Println("~~~ Cleaning plugins cache ~~~")
	os.RemoveAll(plugins.CachePath)
	plugs, err := plugdeps.List(r.App)
	if err != nil && (errx.Unwrap(err) != plugdeps.ErrMissingConfig) {
		return err
	}

	run := genny.WetRunner(context.Background())
	gg, err := install.New(&install.Options{
		App:     r.App,
		Plugins: plugs.List(),
	})

	run.WithGroup(gg)

	fmt.Println("~~~ Reinstalling plugins ~~~")
	return run.Run()
}
