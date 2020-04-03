package env

import (
	"context"
	"os"
	"testing"
)

func Test_BeforeEnv(t *testing.T) {
	beforeTest := &BeforeTester{}

	err := beforeTest.BeforeTest(context.Background(), "", []string{})
	if err != nil {
		t.Fail()
	}

	if os.Getenv("GO_ENV") != "test" {
		t.Error("GO_ENV should be set to be test")
	}
}
