package bzr

import (
	"context"
	"errors"
	"testing"
)

func Test_Bzr_Generalities(t *testing.T) {
	b := Buffalo{}

	if b.Name() != name {
		t.Error("Name should be bzr")
	}

	if b.Description() != description {
		t.Error("description should be something else")
	}
}

func Test_Bzr_BuildVersion(t *testing.T) {

	b := &Buffalo{}

	tcases := []struct {
		name    string
		version string
		err     error

		expectedVersion string
		hasErr          bool
	}{
		{name: "ALL GOOD", version: "abc123", err: nil, expectedVersion: "abc123", hasErr: false},
		{name: "ERROR", version: "", err: errors.New("error loading this thing"), expectedVersion: "abc123", hasErr: true},
	}

	for _, tcase := range tcases {
		testConfig.enabled = true

		t.Run(tcase.name, func(t *testing.T) {
			testConfig.resultError = tcase.err
			testConfig.resultVersion = tcase.version

			result, err := b.BuildVersion(context.Background(), ".")

			if tcase.hasErr && err == nil {
				t.Errorf("Should return err")
			}

			if !tcase.hasErr && err != nil {
				t.Errorf("Should not return err")
			}

			if result != testConfig.resultVersion {
				t.Errorf("Version should be `%+v` but is `%+v`", tcase.expectedVersion, result)
			}
		})
	}
}
