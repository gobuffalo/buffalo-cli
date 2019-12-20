package buildcmd

const mainBuildTmpl = `
package main

import (
	"context"
	"log"
	"os"

  {{if .WithFallthrough -}}
	appcli "{{.Info.Module.Path}}/cli"
	{{end -}}

	{{range $imp := .Imports -}}
	_ "{{$imp}}"
	{{end -}}

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
			Plugger:      cb,
			BuildTime:    {{.BuildTime}},
			BuildVersion: {{.BuildVersion}},
		  {{if .WithFallthrough -}}
			Fallthrough:  appcli.Buffalo,
			{{else -}}
			Fallthrough:  cb.Main,
	    {{end -}}
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
