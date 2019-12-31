package pop

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/pop/v5"
	"github.com/stretchr/testify/require"
)

func Test_Tester_TestArgs(t *testing.T) {
	r := require.New(t)

	tc := &Tester{}

	args, err := tc.TestArgs(nil, "")
	r.NoError(err)
	r.Equal([]string{"-p", "1"}, args)
}

func Test_Tester_BeforeTest(t *testing.T) {
	r := require.New(t)

	ctx := context.Background()

	tc := &Tester{}

	dir, err := ioutil.TempDir("", "")
	r.NoError(err)

	cd := &pop.ConnectionDetails{
		URL: fmt.Sprintf("sqlite://%s", filepath.Join(dir, "test.db")),
	}

	c, err := pop.NewConnection(cd)
	r.NoError(err)
	ctx = context.WithValue(ctx, "tx", c)

	args := []string{}
	err = tc.BeforeTest(ctx, args)
	r.NoError(err)
}
