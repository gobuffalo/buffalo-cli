package main

import (
	"context"
	"log"
	"os"
	"os/exec"

	"github.com/gobuffalo/buffalo-cli/v2/internal/garlic"
)

func Tidy(ctx context.Context) error {
	c := exec.CommandContext(ctx, "go", "mod", "tidy")
	return c.Run()
}

func main() {
	ctx := context.Background()
	defer Tidy(ctx)

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if err := garlic.Run(ctx, pwd, os.Args[1:]); err != nil {
		Tidy(ctx)
		log.Fatal(err)
	}
}
