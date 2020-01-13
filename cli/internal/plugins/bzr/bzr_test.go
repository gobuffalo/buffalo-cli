package bzr

import (
	"bytes"
	"context"
	"errors"
	"testing"
)

type testVersionRunner struct {
	resultError   error
	resultVersion string
}

func (tv *testVersionRunner) ToolAvailable() (bool, error) {
	return true, nil
}

func (tv *testVersionRunner) RunVersionCommand(ctx context.Context, bb *bytes.Buffer) error {
	bb.Write([]byte(tv.resultVersion))
	return tv.resultError
}

func Test_Bzr_Generalities(t *testing.T) {
	b := BzrVersioner{}

	if b.Name() != "bzr" {
		t.Error("Name should be bzr")
	}

	if b.Description() != "Provides bzr related hooks to Buffalo applications." {
		t.Error("description should be something else")
	}
}

func Test_Bzr_BuildVersion(t *testing.T) {

	tcases := []struct {
		name    string
		version string
		err     error

		expectedVersion string
		hasErr          bool
	}{
		{name: "ALL GOOD", version: "abc123", err: nil, expectedVersion: "abc123", hasErr: false},
		{name: "ERROR", version: "", err: errors.New("error loading this thing"), expectedVersion: "", hasErr: true},
	}

	for _, tcase := range tcases {

		t.Run(tcase.name, func(t *testing.T) {

			b := &BzrVersioner{
				versionRunner: &testVersionRunner{
					resultError:   tcase.err,
					resultVersion: tcase.version,
				},
			}

			result, err := b.BuildVersion(context.Background(), ".")

			if tcase.hasErr && err == nil {
				t.Errorf("Should return err")
			}

			if !tcase.hasErr && err != nil {
				t.Errorf("Should not return err")
			}

			if result != tcase.expectedVersion {
				t.Errorf("Version should be `%+v` but is `%+v`", tcase.expectedVersion, result)
			}
		})
	}
}
