package utils

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	app "github.com/upsun/lib-sun"
)

const TEST_CONVERT_PATH = "../../tests/convert/"

func TestBuildTemporyWorkspace(t *testing.T) {
	assert := assert.New(t)

	expected := "/tmp/cli-"
	actual := BuildTemporyWorkspace()
	defer actual.CleanUp()

	assert.Contains(actual.Root, expected)
	assert.DirExists(actual.Root)
	assert.Contains(actual.Data, PATH_DATA)
	assert.Contains(actual.Repo, PATH_REPO)
}

func TestBuildPersistWorkspace(t *testing.T) {
	if os.Getenv("TEST_DEV") == "" {
		t.Skip("Skipping not finished test")
	}

	assert := assert.New(t)

	expected := "./data/"
	actual := BuildPersistWorkspace()

	assert.Contains(actual.Root, expected)
	assert.DirExists(actual.Root)
	assert.Contains(actual.Data, PATH_DATA)
	assert.Contains(actual.Repo, PATH_REPO)
}

func TestMakeDataFolder(t *testing.T) {
	assert := assert.New(t)

	ws := BuildTemporyWorkspace()
	defer ws.CleanUp()

	expected := "test-srv"
	actual := ws.MakeDataFolder("test-srv")

	assert.Contains(actual, ws.Data)
	assert.Contains(actual, expected)
	assert.DirExists(actual)
}

func TestMakeUpsunConfigPath(t *testing.T) {
	assert := assert.New(t)

	base := "/tmp/test-srv"
	attendedPath := ".upsun"
	attendedFile := path.Join(attendedPath, "config.yaml")

	expected := path.Join(base, attendedFile)
	_, actualFile := MakeUpsunConfigPath(base)

	assert.Contains(actualFile, base)
	assert.Contains(actualFile, expected)
	assert.DirExists(strings.Replace(actualFile, "config.yaml", "", 1))
}

func TestCleanUp(t *testing.T) {
	assert := assert.New(t)

	actual := BuildTemporyWorkspace()
	actual.CleanUp()

	assert.NoDirExists(actual.Root)
}

func TestResolvHomePath(t *testing.T) {
	assert := assert.New(t)

	actual := ResolvHomePath("~/test")
	assert.NotContains(actual, "~/")
	assert.True(strings.HasSuffix(actual, "/test"))
}

func TestAbsPath(t *testing.T) {
	assert := assert.New(t)

	expected := "/tmp"
	actual := AbsPath("/tmp")

	assert.Equal(expected, actual)
}

func TestNormalizePath(t *testing.T) {
	assert := assert.New(t)

	actual := NormalizePath("~/test")
	assert.NotContains(actual, "~/")
	assert.True(strings.HasSuffix(actual, "/test"))
}

func TestIsExist(t *testing.T) {
	assert := assert.New(t)

	ws := BuildTemporyWorkspace()
	defer ws.CleanUp()

	actual := IsExist(ws.Root)
	assert.True(actual)

	actual = IsExist(ws.Data)
	assert.True(actual)

	actual = IsExist(ws.Repo)
	assert.True(actual)
}

func TestMakePshConfigPath(t *testing.T) {
	assert := assert.New(t)

	app, psh := MakePshConfigPath("root", "app")
	assert.Equal("root/app/.platform.app.yaml", app)
	assert.Equal("root/.platform", psh)

	app, psh = MakePshConfigPath("root", "")
	assert.Equal("root/.platform.app.yaml", app)
	assert.Equal("root/.platform", psh)
}

func TestIsKeep(t *testing.T) {
	assert := assert.New(t)

	assert.False(IsKeep())

	app.ArgsM.KeepData = "enabled"
	assert.True(IsKeep())
}

func TestCopyFile(t *testing.T) {
	assert := assert.New(t)
	ws := BuildTemporyWorkspace()
	defer ws.CleanUp()

	dst := path.Join(ws.Root, "config.yaml")
	err := CopyFile(TEST_CONVERT_PATH+"config-ref.yaml", dst)
	assert.NoError(err)
	assert.True(IsExist(dst))

	dst = path.Join(ws.Root, "do_not_exist")
	err = CopyFile(TEST_CONVERT_PATH+"do_not_exist", dst)
	assert.Error(err)
	assert.False(IsExist(dst))
}

func TestListDir(t *testing.T) {
	assert := assert.New(t)

	dirs, err := ListDir(TEST_CONVERT_PATH)
	assert.NoError(err)
	assert.Len(dirs, 2)
}
