package main

import (
	"context"
	"log"
	"os"
	"os/exec"

	"github.com/gobuffalo/buffalo-cli/internal/garlic"
)

func Tidy(ctx context.Context) error {
	c := exec.CommandContext(ctx, "go", "mod", "tidy")
	return c.Run()
}

func main() {
	ctx := context.Background()
	defer Tidy(ctx)

	if err := garlic.Run(ctx, os.Args[1:]); err != nil {
		Tidy(ctx)
		log.Fatal(err)
	}
}
