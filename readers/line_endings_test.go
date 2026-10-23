package readers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	entity "github.com/upsun/lib-sun/entity"
	"gopkg.in/yaml.v3"
)

func TestReadServicesCRLF(t *testing.T) {
	const services = `# Database
db:
    type: mariadb:10.6
    # Default disk
    disk: 2048

# Cache
cache:
    type: redis:7.0
`
	dir := t.TempDir()
	lf := filepath.Join(dir, "lf.yaml")
	crlf := filepath.Join(dir, "crlf.yaml")
	require.NoError(t, os.WriteFile(lf, []byte(services), 0o644))
	require.NoError(t, os.WriteFile(crlf, []byte(strings.ReplaceAll(services, "\n", "\r\n")), 0o644))

	var lfModel, crlfModel entity.MetaConfig
	ReadServices(&lfModel, []string{lf})
	ReadServices(&crlfModel, []string{crlf})
	RemoveAllEntry(&lfModel.Services, "disk")
	RemoveAllEntry(&crlfModel.Services, "disk")

	expected, err := yaml.Marshal(&lfModel)
	require.NoError(t, err)
	actual, err := yaml.Marshal(&crlfModel)
	require.NoError(t, err)
	assert.Equal(t, string(expected), string(actual))
}

func TestReadPlatformsSourceRoot(t *testing.T) {
	workspace := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(workspace, "sub"), 0o755))
	rootApp := filepath.Join(workspace, entity.PSH_PLATFORM)
	subApp := filepath.Join(workspace, "sub", entity.PSH_PLATFORM)
	require.NoError(t, os.WriteFile(rootApp, []byte("name: app\ntype: php:8.2\n"), 0o644))
	require.NoError(t, os.WriteFile(subApp, []byte("name: sub\ntype: php:8.2\n"), 0o644))

	var metamodel entity.MetaConfig
	ReadPlatforms(&metamodel, []string{rootApp, subApp}, workspace)

	out, err := yaml.Marshal(&metamodel)
	require.NoError(t, err)
	assert.Contains(t, string(out), "        source:\n            root: /\n")
	assert.Contains(t, string(out), "        source:\n            root: /sub/\n")
}
