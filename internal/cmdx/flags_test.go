package cmdx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FlagSet_Parse(t *testing.T) {
	r := require.New(t)

	flags := NewFlagSet("")

	var known bool
	var good bool
	flags.BoolVarP(&known, "known", "k", false, "")
	flags.BoolVarP(&good, "good", "g", false, "")

	args := []string{"--known", "--unknown", "--good", "--bad"}

	err := flags.Parse(args)
	r.Error(err)

	r.True(known)
	r.True(good)

	u, ok := err.(Unknowns)
	r.True(ok)
	r.Len(u.Unknowns(), 2)
}
