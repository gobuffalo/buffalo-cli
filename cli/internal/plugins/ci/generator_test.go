package ci

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Generator(t *testing.T) {
	r := require.New(t)

	dir, err := ioutil.TempDir("", "tmp")
	r.NoError(err)
	defer os.RemoveAll(dir)

	generator := &Generator{}
	ctx := context.Background()

	err = generator.Newapp(ctx, dir, "app", []string{})
	r.NoError(err)

	b, err := ioutil.ReadFile(filepath.Join(dir, ".github/workflows/test.yml"))
	r.NoError(err)

	r.Contains(string(b), "name: Tests")
	r.Contains(string(b), "${{ runner.os }}-go")
	r.Contains(string(b), `TEST_DATABASE_URL: "postgres://postgres:postgres@127.0.0.1:${{ job.services.postgres.ports[5432] }}/app_test?sslmode=disable"`)
}
