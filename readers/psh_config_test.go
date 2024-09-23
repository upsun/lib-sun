package readers

import (
	"path"
	"testing"

	assert "github.com/stretchr/testify/assert"
	entity "github.com/upsun/lib-sun/entity"
	"gopkg.in/yaml.v3"
)

const TEST_CONVERT_PATH = "../tests/project-psh"
const TEST_PSH_PATH = ".platform/"

func TestReadServices(t *testing.T) {
	assert := assert.New(t)

	var metamodel entity.MetaConfig

	ReadServices(&metamodel, []string{path.Join(TEST_CONVERT_PATH, TEST_PSH_PATH, entity.PSH_SERVICE)})

	assert.Equal(yaml.MappingNode, metamodel.Services.Kind)
	assert.Equal(TAG_MAP, metamodel.Services.Tag)
	assert.Len(metamodel.Services.Content, 10)
}

func TestReadApplications(t *testing.T) {
	assert := assert.New(t)

	var metamodel entity.MetaConfig

	ReadPlatforms(&metamodel, []string{path.Join(TEST_CONVERT_PATH, entity.PSH_PLATFORM)}, TEST_CONVERT_PATH)

	assert.Equal(yaml.MappingNode, metamodel.Applications.Kind)
	assert.Equal(TAG_MAP, metamodel.Applications.Tag)
	assert.Len(metamodel.Applications.Content, 2)
}

func TestReadRoutes(t *testing.T) {
	assert := assert.New(t)

	var metamodel entity.MetaConfig

	ReadRoutes(&metamodel, []string{path.Join(TEST_CONVERT_PATH, TEST_PSH_PATH, entity.PSH_ROUTE)})

	assert.Equal(yaml.MappingNode, metamodel.Routes.Kind)
	assert.Equal(TAG_MAP, metamodel.Routes.Tag)
	assert.Len(metamodel.Routes.Content, 4)
}
