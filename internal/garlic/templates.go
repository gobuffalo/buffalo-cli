package garlic

const cliBuffalo = `
package cli

import (
	"context"
	"fmt"

	"github.com/gobuffalo/buffalo-cli/cli"
)

func Buffalo(ctx context.Context, args []string) error {
	fmt.Println("~~~~ Using {{.Name}}/cli.Buffalo ~~~")
	buffalo, err := cli.New()
	if err != nil {
		return err
	}

	// buffalo.Plugins = append(buffalo.Plugins,
	// 	your plugins here!
	// )
	return buffalo.Main(ctx, args)
}
`
