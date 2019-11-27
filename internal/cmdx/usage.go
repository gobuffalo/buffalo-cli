package cmdx

import (
	"context"
	"flag"
	"fmt"
)

func Usage(ctx context.Context, flags *flag.FlagSet) {
	stderr := Stderr(ctx)
	flags.SetOutput(stderr)
	flags.Usage = func() {
		fmt.Fprintf(stderr, "Usage of %s:\n", flags.Name())
		flags.PrintDefaults()
	}
}
