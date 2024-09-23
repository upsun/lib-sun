package detector

import (
	"path"
	"testing"

	assert "github.com/stretchr/testify/assert"
	entity "github.com/upsun/lib-sun/entity"
)

const TEST_PATH = "../tests/project-psh"

func TestFind(t *testing.T) {
	assert := assert.New(t)

	want := make(map[string][]string)
	want[entity.PSH_ROUTE] = []string{path.Join(TEST_PATH, PSH_CONFIG_PATH, entity.PSH_ROUTE)}
	want[entity.PSH_SERVICE] = []string{path.Join(TEST_PATH, PSH_CONFIG_PATH, entity.PSH_SERVICE)}
	want[entity.PSH_APPLICATION] = []string{path.Join(TEST_PATH, PSH_CONFIG_PATH, entity.PSH_APPLICATION)}
	want[entity.PSH_PLATFORM] = []string{path.Join(TEST_PATH, entity.PSH_PLATFORM)}

	got, _ := FindConfig(TEST_PATH)

	for key, values := range want {
		assert.GreaterOrEqual(len(values), 1)

		if len(values) != 0 {
			assert.Equal(values[0], got[key][0])
		}
	}

}
