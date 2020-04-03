package test

import (
	"context"
	"os"
	"testing"
)

func Test_GoEnv(t *testing.T) {
	goEnv := &GoEnv{}

	os.Setenv("GO_ENV", "Antonio's Home")

	err := goEnv.BeforeTest(context.Background(), "", []string{})
	if err != nil {
		t.Error("error setting GO_ENV")
	}

	if os.Getenv("GO_ENV") != "test" {
		t.Error("GO_ENV should be set to be test")
	}

	err = goEnv.AfterTest(context.Background(), "", []string{}, nil)
	if err != nil {
		t.Error("error setting GO_ENV")
	}

	if os.Getenv("GO_ENV") != "Antonio's Home" {
		t.Error("GO_ENV should be set to be blank")
	}
}
