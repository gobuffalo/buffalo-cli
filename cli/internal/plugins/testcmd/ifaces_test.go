package testcmd

import (
	"context"
	"os/exec"
)

var _ Tester = &tester{}

type tester struct {
	name string
	args []string
	err  error
}

func (b *tester) Name() string {
	if len(b.name) == 0 {
		return "tester"
	}
	return b.name
}

func (b *tester) Test(ctx context.Context, args []string) error {
	b.args = args
	return b.err
}

var _ BeforeTester = &beforeTester{}

type beforeTester struct {
	name string
	args []string
	err  error
}

func (b *beforeTester) Name() string {
	if len(b.name) == 0 {
		return "beforeTester"
	}
	return b.name
}

func (b *beforeTester) BeforeTest(ctx context.Context, args []string) error {
	b.args = args
	return b.err
}

var _ AfterTester = &afterTester{}

type afterTester struct {
	name string
	args []string
	err  error
}

func (b *afterTester) Name() string {
	if len(b.name) == 0 {
		return "afterTester"
	}
	return b.name
}

func (b *afterTester) AfterTest(ctx context.Context, args []string, err error) error {
	b.args = args
	b.err = err
	return err
}

var _ Tagger = &tagger{}

type tagger struct {
	root string
	tags []string
	err  error
}

func (b *tagger) TestTags(ctx context.Context, root string) ([]string, error) {
	b.root = root
	return b.tags, b.err
}

var _ Runner = &bladeRunner{}

type bladeRunner struct {
	cmd *exec.Cmd
	err error
}

func (bladeRunner) Name() string {
	return "blade"
}

func (b *bladeRunner) RunTests(ctx context.Context, cmd *exec.Cmd) error {
	b.cmd = cmd
	return b.err
}
