package test

import (
	"context"
	"os"
	"testing"
)

func Test_GoEnv_BeforeTest(t *testing.T) {
	goEnv := &GoEnv{}

	err := goEnv.BeforeTest(context.Background(), "", []string{})
	if err != nil {
		t.Error("error setting GO_ENV")
	}

	if os.Getenv("GO_ENV") != "test" {
		t.Error("GO_ENV should be set to be test")
	}
}

func Test_GoEnv_AfterTest(t *testing.T) {
	goEnv := &GoEnv{}

	os.Setenv("GO_ENV", "")

	err := goEnv.AfterTest(context.Background(), "", []string{}, nil)
	if err != nil {
		t.Error("error setting GO_ENV")
	}

	if os.Getenv("GO_ENV") != "" {
		t.Error("GO_ENV should be set to be blank")
	}
}
