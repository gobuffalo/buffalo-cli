package cmdx

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/spf13/pflag"
)

func NewFlagSet(name string, w io.Writer) *pflag.FlagSet {
	flags := pflag.NewFlagSet(name, pflag.ContinueOnError)
	flags.SetOutput(ioutil.Discard)
	flags.Usage = func() {
		fmt.Fprintf(w, "Usage of %s:\n", name)
		flags.SetOutput(w)
		flags.PrintDefaults()
	}
	return flags
}
