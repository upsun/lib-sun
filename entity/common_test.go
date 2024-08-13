package entity

import (
	"testing"

	assert "github.com/stretchr/testify/assert"
)

const (
	PROVIDER = "provider"
	ID       = "id"
	ENV      = "env"
	KEY      = "test_key"
)

func TestMakeProjectContext(t *testing.T) {
	assert := assert.New(t)

	actual := MakeProjectContext(PROVIDER, ID, ENV)

	assert.Equal(PROVIDER, actual.Provider)
	assert.Equal(ID, actual.ID)
	assert.Equal(ENV, actual.DefaultEnv)
	assert.NotNil(actual.Variables)
	assert.NotNil(actual.Users)
	assert.NotNil(actual.Access)
	assert.NotNil(actual.Services)
	assert.NotNil(actual.Mounts)
}

func TestProjectGlobal_CopyProjectBase(t *testing.T) {
	assert := assert.New(t)

	base := MakeProjectContext(PROVIDER, ID, ENV)
	base.Name = "name"
	base.Description = "description"
	base.Timezone = "timezone"
	base.Region = "region"
	base.DefaultDom = "domaine"

	actual := MakeProjectContext("", "", "")
	actual.CopyProjectBase(base)

	assert.Equal(base.Name, actual.Name)
	assert.Equal(base.Description, actual.Description)
	assert.Equal(base.DefaultEnv, actual.DefaultEnv)
	assert.Equal(base.Timezone, actual.Timezone)
	assert.Equal(base.Region, actual.Region)
	assert.Equal(base.DefaultDom, actual.DefaultDom)
}

func TestProjectGlobal_CopyVariables(t *testing.T) {
	assert := assert.New(t)

	base := MakeProjectContext(PROVIDER, ID, ENV)
	actual := MakeProjectContext("", "", "")
	actual.CopyVariables(base)
	assert.Len(actual.Variables, 0)

	base.Variables[KEY] = ProjectVariable{}
	actual.CopyVariables(base)
	assert.Len(actual.Variables, 1)
	assert.NotNil(actual.Variables[KEY])
}

func TestProjectGlobal_CopyUsers(t *testing.T) {
	assert := assert.New(t)

	base := MakeProjectContext(PROVIDER, ID, ENV)
	actual := MakeProjectContext("", "", "")
	actual.CopyUsers(base)
	assert.Len(actual.Users, 0)

	base.Users[KEY] = ProjectUser{}
	actual.CopyUsers(base)
	assert.Len(actual.Users, 1)
	assert.NotNil(actual.Users[KEY])
}

func TestProjectGlobal_CopyAccess(t *testing.T) {
	assert := assert.New(t)

	base := MakeProjectContext(PROVIDER, ID, ENV)
	actual := MakeProjectContext("", "", "")
	actual.CopyAccess(base)
	assert.Len(actual.Access, 0)

	base.Access[KEY] = ProjectAccess{}
	actual.CopyAccess(base)
	assert.Len(actual.Access, 1)
	assert.NotNil(actual.Access[KEY])
}

func TestProjectGlobal_CopyServices(t *testing.T) {
	assert := assert.New(t)

	base := MakeProjectContext(PROVIDER, ID, ENV)
	actual := MakeProjectContext("", "", "")
	actual.CopyServices(base)
	assert.Len(actual.Services, 0)

	base.Services[KEY] = EnvService{}
	actual.CopyServices(base)
	assert.Len(actual.Services, 1)
	assert.NotNil(actual.Services[KEY])
}

func TestProjectGlobal_CopyMounts(t *testing.T) {
	assert := assert.New(t)

	base := MakeProjectContext(PROVIDER, ID, ENV)
	actual := MakeProjectContext("", "", "")
	actual.CopyMounts(base)
	assert.Len(actual.Mounts, 0)

	base.Mounts[KEY] = EnvMount{}
	actual.CopyMounts(base)
	assert.Len(actual.Mounts, 1)
	assert.NotNil(actual.Mounts[KEY])
}

func TestProjectGlobal_Copy(t *testing.T) {
	assert := assert.New(t)

	base := MakeProjectContext(PROVIDER, ID, ENV)
	base.Variables[KEY] = ProjectVariable{}
	base.Users[KEY] = ProjectUser{}
	base.Access[KEY] = ProjectAccess{}
	base.Services[KEY] = EnvService{}
	base.Mounts[KEY] = EnvMount{}

	actual := MakeProjectContext(PROVIDER, ID, ENV)
	actual.Copy(base)

	assert.EqualValues(base, actual)
}
