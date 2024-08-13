package readers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	entity "github.com/upsun/lib-sun/entity"
	"gopkg.in/yaml.v3"
)

const TEST_CONVERT_PATH = "../../tests/convert/"

func TestReadServices(t *testing.T) {
	assert := assert.New(t)

	var metamodel entity.MetaConfig

	ReadServices(&metamodel, []string{TEST_CONVERT_PATH + ".platform/services.yaml"})

	assert.Equal(yaml.MappingNode, metamodel.Services.Kind)
	assert.Equal(TAG_MAP, metamodel.Services.Tag)
	assert.Len(metamodel.Services.Content, 10)
}

func TestReadApplications(t *testing.T) {
	assert := assert.New(t)

	var metamodel entity.MetaConfig

	ReadPlatforms(&metamodel, []string{TEST_CONVERT_PATH + ".platform.app.yaml"}, TEST_CONVERT_PATH)

	assert.Equal(yaml.MappingNode, metamodel.Applications.Kind)
	assert.Equal(TAG_MAP, metamodel.Applications.Tag)
	assert.Len(metamodel.Applications.Content, 2)
}

func TestReadRoutes(t *testing.T) {
	assert := assert.New(t)

	var metamodel entity.MetaConfig

	ReadRoutes(&metamodel, []string{TEST_CONVERT_PATH + ".platform/routes.yaml"})

	assert.Equal(yaml.MappingNode, metamodel.Routes.Kind)
	assert.Equal(TAG_MAP, metamodel.Routes.Tag)
	assert.Len(metamodel.Routes.Content, 4)
}
