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

	tcases := []struct {
		provider string
		filepath string
		contents []string
	}{
		{
			provider: "",
			filepath: filepath.Join(dir, ".github/workflows/test.yml"),
			contents: []string{
				"name: Tests",
				"${{ runner.os }}-go",
				`TEST_DATABASE_URL: "postgres://postgres:postgres@127.0.0.1:${{ job.services.postgres.ports[5432] }}/app_test?sslmode=disable"`,
			},
		},
		{
			provider: "github",
			filepath: filepath.Join(dir, ".github/workflows/test.yml"),
			contents: []string{
				"name: Tests",
				"${{ runner.os }}-go",
				`TEST_DATABASE_URL: "postgres://postgres:postgres@127.0.0.1:${{ job.services.postgres.ports[5432] }}/app_test?sslmode=disable"`,
			},
		},

		{
			provider: "travis",
			filepath: filepath.Join(dir, ".travis.yml"),
			contents: []string{
				"language: go",
				"- psql -c 'create database app_test;' -U postgres",
			},
		},

		{
			provider: "gitlab",
			filepath: filepath.Join(dir, ".gitlab-ci.yml"),
			contents: []string{
				"- apt-get update && apt-get install -y postgresql-client",
				"before_script:",
			},
		},

		{
			provider: "circleci",
			filepath: filepath.Join(dir, ".circleci", "config.yml"),
			contents: []string{
				"version: 2",
				"jobs:",
				"- image: circleci/postgres:9.6-alpine",
			},
		},
	}

	for index, tcase := range tcases {
		generator.provider = tcase.provider

		err = generator.Newapp(ctx, dir, "app", []string{})
		r.NoError(err)

		b, err := ioutil.ReadFile(tcase.filepath)
		r.NoError(err)

		for _, content := range tcase.contents {
			r.Contains(string(b), content, "Should contain %v - %v", content, index)
		}
	}

}
