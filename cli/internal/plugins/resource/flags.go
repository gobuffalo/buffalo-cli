package resource

import "github.com/spf13/pflag"

func (g *Generator) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(g.Name(), pflag.ContinueOnError)
	flags.BoolVarP(&g.help, "help", "h", false, "print this help")
	return flags
}
