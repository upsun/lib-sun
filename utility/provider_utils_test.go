package utils

import (
	"os"
	"os/exec"
	"testing"

	assert "github.com/stretchr/testify/assert"
	entity "github.com/upsun/lib-upsun/entity"
)

const PSH_PRJ_ID = "a6gx2dq4x235u"

func TestClean(t *testing.T) {
	assert := assert.New(t)

	var actual PshResult
	actual.Title = "title"
	actual.Code = -1
	actual.Status = "status"
	actual.Message = "message"

	actual.Clean()

	assert.Empty(actual.Status)
	assert.Equal(0, actual.Code)
	assert.Empty(actual.Title)
	assert.Empty(actual.Message)
}

func TestCallExePath(t *testing.T) {
	assert := assert.New(t)

	payload := []string{"OK"}
	stdout, stderr, err := CallExePath("echo", "", "/tmp", payload...)

	assert.Equal("OK\n", string(stdout[:]))
	assert.Empty(stderr)
	assert.NoError(err)
}

func TestCallExe(t *testing.T) {
	assert := assert.New(t)

	payload := []string{"OK"}
	stdout, stderr, err := CallExe("echo", "", payload...)

	assert.Equal("OK\n", string(stdout[:]))
	assert.Empty(stderr)
	assert.NoError(err)
}

func TestCallCLIPath(t *testing.T) {
	assert := assert.New(t)

	payload := []string{}
	prj := entity.MakeProjectContext(entity.PSH_PROVIDER, "", "")
	stdout, err := CallCLIPath(prj, "version", "/tmp", payload...)

	assert.Contains(string(stdout[:]), "Platform.sh CLI")
	assert.NoError(err)
}

func TestCallCLI(t *testing.T) {
	assert := assert.New(t)

	payload := []string{}
	prj := entity.MakeProjectContext(entity.PSH_PROVIDER, "", "")
	stdout, err := CallCLI(prj, "version", payload...)

	assert.Contains(string(stdout[:]), "Platform.sh CLI")
	assert.NoError(err)
}

func TestCallCLIString(t *testing.T) {
	assert := assert.New(t)

	payload := []string{}
	prj := entity.MakeProjectContext(entity.PSH_PROVIDER, "", "")
	stdout, err := CallCLIString(prj, "version", payload...)

	assert.Contains(stdout, "Platform.sh CLI")
	assert.NoError(err)
}

func TestCallAPI(t *testing.T) {
	assert := assert.New(t)

	payload := []string{"/system"}
	prj := entity.MakeProjectContext(entity.PSH_PROVIDER, PSH_PRJ_ID, "master")
	pshResult := CallAPI(prj, payload...)

	// TODO make real test with mock.
	assert.Equal(0, pshResult.Code)
}

func TestHackAlias(t *testing.T) {
	assert := assert.New(t)

	for _, provider := range PROVIDER_SUPPORTED {
		cmd, env := HackAlias(provider)
		assert.Equal(entity.PSH_PROVIDER, cmd)

		if provider == entity.PSH_PROVIDER {
			assert.Empty(env)
		} else {
			assert.Contains(env, "CLI_CONFIG_FILE=")
		}
	}
}

func TestIsAuthenticated(t *testing.T) {
	if os.Getenv("TEST_DEV") == "" {
		t.Skip("Skipping not finished test")
	}
	assert := assert.New(t)

	prj := entity.MakeProjectContext(entity.PSH_PROVIDER, "", "")
	auth := IsAuthenticated(prj)

	assert.True(auth)
}

func TestHasSufficientRights(t *testing.T) {
	if os.Getenv("TEST_DEV") == "" {
		t.Skip("Skipping not finished test")
	}
}

func TestProviderCheck(t *testing.T) {
	assert := assert.New(t)

	if os.Getenv("BE_CRASHER") == "1" {
		ProviderCheck(entity.PSH_PROVIDER)
		return
	}

	if os.Getenv("BE_CRASHER") == "2" {
		ProviderCheck("foo")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestProviderCheck")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	assert.NoError(err)

	cmd = exec.Command(os.Args[0], "-test.run=TestProviderCheck")
	cmd.Env = append(os.Environ(), "BE_CRASHER=2")
	err = cmd.Run()
	assert.Error(err)
}

func TestCliJsonString(t *testing.T) {
	assert := assert.New(t)

	actual := CliJsonString("key", "value")
	assert.Equal("\"key\":\"value\"", actual)
}

func TestCliJsonBoolean(t *testing.T) {
	assert := assert.New(t)

	actual := CliJsonBoolean("key", "true")
	assert.Equal("\"key\":true", actual)
}

func TestCliJson(t *testing.T) {
	assert := assert.New(t)

	actual := CliJson("key1", "key2", "key3")
	assert.Equal("{key1,key2,key3}", actual)
}
