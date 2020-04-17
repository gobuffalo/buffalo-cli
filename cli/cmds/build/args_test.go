package build

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_AppendArg(t *testing.T) {
	t.Parallel()

	table := []struct {
		arg  string
		exp  []string
		in   []string
		name string
	}{
		{name: "-tags", in: []string{"-tags", "foo"}, arg: "bar", exp: []string{"-tags", "foo bar"}},
	}

	for _, tt := range table {
		t.Run(tt.name, func(st *testing.T) {
			r := require.New(st)

			act := AppendArg(tt.in, tt.name, tt.arg)
			r.Equal(tt.exp, act)
		})
	}

}
