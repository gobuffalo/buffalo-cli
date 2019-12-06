package cmdx

import (
	"io/ioutil"

	"github.com/spf13/pflag"
)

// NewFlagSet ...
func NewFlagSet(name string) *pflag.FlagSet {
	flags := pflag.NewFlagSet(name, pflag.ContinueOnError)
	flags.SetOutput(ioutil.Discard)
	return flags
}
