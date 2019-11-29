package cli

import (
	"flag"
	"io/ioutil"
)

func NewFlagSet(name string) *flag.FlagSet {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(ioutil.Discard)
	return flags
}
