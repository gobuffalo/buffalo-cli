package portal

import (
	"context"
	"fmt"
	"net/url"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/develop"
)

var _ develop.Developer = &Developer{}

type Developer struct{}

func (Developer) PluginName() string {
	return "portal/developer"
}

func (d *Developer) Develop(ctx context.Context, root string, args []string) error {
	u, err := url.Parse("http://0.0.0.0:3000")
	if err != nil {
		return err
	}
	fmt.Println(">>>TODO cli/internal/plugins/portal/developer.go:19: u ", u)
	return nil
}
