package utils

import (
	"os"
	"testing"

	flag "github.com/spf13/pflag"
	assert "github.com/stretchr/testify/assert"
)

func TestDisclaimer(t *testing.T) {
	Disclaimer("Unit Test")
}

func TestIsFlagPassed(t *testing.T) {
	assert := assert.New(t)
	os.Args[1] = "--test=bar"
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	var test string
	flag.StringVarP(&test, "test", "t", "default", "Test")

	assert.False(IsFlagPassed("verbose"))
	assert.Equal("default", test)

	flag.Parse()

	assert.Equal("bar", test)
	assert.True(IsFlagPassed("test"))
}

func TestRequireFlag(t *testing.T) {
	assert := assert.New(t)
	os.Args[1] = "--test=bar"
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	var test string
	flag.StringVarP(&test, "test", "t", "default", "Test")
	flag.Parse()

	assert.Equal("bar", RequireFlag("test", "No question on unit test", test, false))
	assert.True(IsFlagPassed("test"))

	//TODO test case interactive
}
