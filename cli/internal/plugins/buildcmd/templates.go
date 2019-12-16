package buildcmd

const mainBuildTmpl = `
package main

import (
	"context"
	"log"
	"os"

	appcli "coke/cli"

	// import all application packages
	_ "coke/actions"
	_ "coke/grifts"
	_ "coke/models"

	// etc...

	"github.com/gobuffalo/buffalo-cli/built"
	"github.com/gobuffalo/buffalo-cli/cli"
)

func main() {
	err := func() error {
		cb, err := cli.New()
		if err != nil {
			return err
		}

		b := &built.App{
			Buffalo:      cb,
			BuildTime:    {{.BuildTime}},
			BuildVersion: {{.BuildVersion}},
			Fallthrough:  appcli.Buffalo,
			OriginalMain: originalMain,
		}

		ctx := context.Background()
		return b.Main(ctx, os.Args[1:])
	}()

	if err != nil {
		log.Fatal(err)
	}
}
`
