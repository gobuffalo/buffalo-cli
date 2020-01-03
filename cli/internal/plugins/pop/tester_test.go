// +build sqlite

package pop

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	dir := filepath.Join("internal", "testdata", "temp")
	i := m.Run()
	os.RemoveAll(dir)
	os.Exit(i)
}

func Test_Tester_TestArgs(t *testing.T) {
	r := require.New(t)

	tc := &Tester{}

	args, err := tc.TestArgs(nil, "")
	r.NoError(err)
	r.Equal([]string{"-p", "1"}, args)
}

func Test_Tester_BeforeTest_widgets_migrations(t *testing.T) {
	r := require.New(t)

	tc := &Tester{}

	ref, err := testerRef()
	r.NoError(err)

	mf := filepath.Join(ref.Dir, "migrations", "1_widgets.up.fizz")
	r.NoError(writeFile(mf, dbWidgetsMigration))

	tc.WithHereInfo(ref.Info)

	args := []string{}

	err = tc.BeforeTest(ref.Context(), args)
	r.NoError(err)

	tx := ref.TX
	count, err := tx.Count("widgets")
	r.NoError(err)
	r.Equal(0, count)
}

func Test_Tester_BeforeTest_force_migrations(t *testing.T) {
	r := require.New(t)

	tc := &Tester{}

	ref, err := testerRef()
	r.NoError(err)

	r.NoError(writeSchema(ref.Info, dbEmptySchema))

	mf := filepath.Join(ref.Dir, "migrations", "1_widgets.up.fizz")
	r.NoError(writeFile(mf, dbWidgetsMigration))

	tc.WithHereInfo(ref.Info)

	args := []string{"--force-migrations"}

	err = tc.BeforeTest(ref.Context(), args)
	r.NoError(err)

	tx := ref.TX
	count, err := tx.Count("widgets")
	r.NoError(err)
	r.Equal(0, count)
}

func Test_Tester_BeforeTest_widgets_schema(t *testing.T) {
	r := require.New(t)

	tc := &Tester{}

	ref, err := testerRef()
	r.NoError(err)

	r.NoError(writeSchema(ref.Info, dbWidgetsSchema))

	tc.WithHereInfo(ref.Info)

	args := []string{}

	err = tc.BeforeTest(ref.Context(), args)
	r.NoError(err)

	tx := ref.TX
	count, err := tx.Count(tx.MigrationTableName())
	r.NoError(err)
	r.Equal(0, count)
}

func Test_Tester_BeforeTest_empty_schema(t *testing.T) {
	r := require.New(t)

	tc := &Tester{}

	ref, err := testerRef()
	r.NoError(err)

	r.NoError(writeSchema(ref.Info, dbEmptySchema))

	tc.WithHereInfo(ref.Info)

	args := []string{}

	err = tc.BeforeTest(ref.Context(), args)
	r.NoError(err)
}

func Test_Tester_BeforeTest_no_schema(t *testing.T) {
	r := require.New(t)

	tc := &Tester{}

	ref, err := testerRef()
	r.NoError(err)

	tc.WithHereInfo(ref.Info)

	args := []string{}
	err = tc.BeforeTest(ref.Context(), args)
	r.NoError(err)
}
