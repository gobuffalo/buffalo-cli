package cmdx

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/pflag"
)

type FlagSet struct {
	*pflag.FlagSet
	args []string
}

func (f *FlagSet) Args() []string {
	return f.args
}

func (flags *FlagSet) Parse(args []string) error {
	flags.args = args
	var unknowns []error
	for _, a := range args {
		if err := flags.FlagSet.Parse([]string{a}); err != nil {
			unknowns = append(unknowns, err)
		}
	}
	if len(unknowns) > 0 {
		return unknownFlagsError{unknowns: unknowns}
	}
	return nil
}

type unknownFlagsError struct {
	unknowns []error
}

func (u unknownFlagsError) Error() string {
	return fmt.Sprintf("unknown flags: %v", u.unknowns)
}

func (u unknownFlagsError) Unknowns() []error {
	return u.unknowns
}

type Unknowns interface {
	Unknowns() []error
}

// NewFlagSet ...
func NewFlagSet(n string) *FlagSet {
	flags := pflag.NewFlagSet(n, pflag.ContinueOnError)
	flags.SetOutput(ioutil.Discard)
	return &FlagSet{
		FlagSet: flags,
	}
}
