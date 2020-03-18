package build

import (
	"strings"
)

func AppendArg(args []string, name string, arg ...string) []string {
	for i, a := range args {
		if a != name {
			continue
		}
		if len(args) <= i {
			return args
		}
		v := args[i+1]
		x := []string{v}
		x = append(x, arg...)
		args[i+1] = strings.Join(x, " ")
		return args
	}
	args = append(args, name, strings.Join(arg, " "))
	return args
}
