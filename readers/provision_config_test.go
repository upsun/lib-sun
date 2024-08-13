package readers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	entity "github.com/upsun/lib-sun/entity"
	utils "github.com/upsun/lib-sun/utility"
)

func TestReadProvisionFile(t *testing.T) {
	assert := assert.New(t)
	ws := utils.BuildTemporyWorkspace()
	defer ws.CleanUp()

	var provision entity.ProvisionGlobal
	configPath := "../../tests/provision"

	projectDstContext := entity.MakeProjectContext(entity.UPS_PROVIDER, "", "master")
	ReadProvisionFile(&projectDstContext, &provision, configPath, ws)

	// Provision entity
	assert.Equal("Bench-Static-file-0", provision.Name)
	assert.Equal("", provision.Applications["app"].Mainstream.URL)
	assert.Equal("master", provision.Applications["app"].Mainstream.Version)
	assert.Equal("TEST", provision.Variables["env:provision"])
	assert.Len(provision.Applications["app"].Files, 0)

	// Project entity
	assert.Equal("Bench-Static-file-0", projectDstContext.Name)
	assert.Equal("TEST", projectDstContext.Variables["env:provision"].Value)
}
