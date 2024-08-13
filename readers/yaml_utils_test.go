package readers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

const KEY_SIZE = "size"

var node_root yaml.Node

func setupTest(_ testing.TB) func(tb testing.TB) {
	fmt.Println(">>> SETUP test")

	var (
		yaml_string = `version: 1
# Test comment
services:
    sqldb:
        # (https://docs.platform.sh/configuration/services/mysql.html#supported-versions)
        type: mysql:10.5
        disk: 1024
        size: M
`
	)

	err := yaml.Unmarshal([]byte(yaml_string), &node_root)
	if err != nil {
		fmt.Printf("Error parsing YAML: %s\n", err)
	}

	return func(tb testing.TB) {
		fmt.Println("<<< TEARDOWN test")
	}
}

func TestNewKeyValueNodes(t *testing.T) {
	assert := assert.New(t)

	got := NewKeyValueNodes()

	assert.NotNil(got.KeyNode, "KeyNode not initialize !")
	assert.NotNil(got.ValueNode, "ValueNode not initialize !")
}

func TestCreateEntryEmpty(t *testing.T) {
	assert := assert.New(t)

	gotValue := "test-CreateEntryEmpty"
	got := CreateEntryEmpty(gotValue)

	if assert.NotNil(got.KeyNode, "KeyNode not initialize !") {
		assert.Equal(yaml.ScalarNode, got.KeyNode.Kind, "KeyNode not good king.")
		assert.Equal(TAG_STRING, got.KeyNode.Tag, "KeyNode not good tag.")
		assert.Equal(gotValue, got.KeyNode.Value, "KeyNode not good value.")
	}

	assert.NotNil(got.ValueNode, "KeyValue not initialize !")
}

func TestCreateEntryString(t *testing.T) {
	assert := assert.New(t)

	gotKey := "test-CreateEntryString-key"
	gotValue := "test-CreateEntryString-value"
	got := CreateEntryString(gotKey, gotValue)

	assert.Equal(gotKey, got.KeyNode.Value, "KeyNode not good value.")

	if assert.NotNil(got.ValueNode, "ValueNode not initialize !") {
		assert.Equal(yaml.ScalarNode, got.ValueNode.Kind, "ValueNode not good king.")
		assert.Equal(TAG_STRING, got.ValueNode.Tag, "ValueNode not good tag.")
		assert.Equal(gotValue, got.ValueNode.Value, "ValueNode not good value.")
	}
}

func TestCreateEntryEntity(t *testing.T) {
	assert := assert.New(t)

	gotKey := "test-CreateEntryEntity-key"
	gotValueValue := "test-CreateEntryEntity-value"
	gotValue := CreateEntryEmpty(gotValueValue)
	got := CreateEntryEntity(gotKey, gotValue)

	assert.Equal(gotKey, got.KeyNode.Value, "KeyNode not good value.")

	if assert.NotNil(got.ValueNode, "ValueNode not initialize !") {
		assert.Equal(yaml.MappingNode, got.ValueNode.Kind, "ValueNode not good king.")
		assert.Equal(TAG_MAP, got.ValueNode.Tag, "ValueNode not good tag.")
		assert.Len(got.ValueNode.Content, 2, "KeyValue not good value")
	}
}

func TestAddEntry(t *testing.T) {
	assert := assert.New(t)

	gotKey := "test-AddEntry-key"
	gotValue := "test-AddEntry-value"
	entrie := CreateEntryString(gotKey, gotValue)
	node := &yaml.Node{}

	AddEntry(node, entrie)

	assert.Equal(gotKey, node.Content[0].Value)
	assert.Equal(gotValue, node.Content[1].Value)
}

func TestFindEntry(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)
	assert := assert.New(t)

	// Check key
	want := "sqldb"
	got := FindEntry(&node_root, want)
	assert.Equal(want, got.KeyNode.Value)

	// Check value
	got = FindEntry(&node_root, KEY_SIZE)
	want = "M"
	assert.Equal(want, got.ValueNode.Value)
}

func TestRemoveEntry(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)
	assert := assert.New(t)

	// Check value
	got := FindEntry(&node_root, KEY_SIZE)
	want := "M"
	assert.Equal(want, got.ValueNode.Value)

	parent := RemoveEntry(&node_root, KEY_SIZE)
	got = FindEntry(&node_root, KEY_SIZE)

	assert.Nil(got.ValueNode, "Node not Remove")
	assert.NotNil(parent, "Node not Remove.")

	parent = RemoveEntry(&node_root, "unknow")
	assert.Nil(parent, "Node found !")

	content, err := yaml.Marshal(&node_root)
	assert.NoError(err)
	assert.NotEmpty(content)
}

func TestRemoveAllEntry(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)
	assert := assert.New(t)

	// Check value
	got := FindEntry(&node_root, KEY_SIZE)
	want := "M"
	assert.Equal(want, got.ValueNode.Value)

	RemoveAllEntry(&node_root, KEY_SIZE)
	got = FindEntry(&node_root, KEY_SIZE)

	assert.Nil(got.ValueNode, "Node not Remove")

	RemoveAllEntry(&node_root, "unknow")

	content, err := yaml.Marshal(&node_root)
	assert.NoError(err)
	assert.NotEmpty(content)
}

func TestReplaceEntry(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)
	assert := assert.New(t)

	node := ReplaceEntry(&node_root, "mysql:10.5", "mariadb:10.5")
	assert.NotNil(node, "Node not replace")

	content, err := yaml.Marshal(&node_root)
	assert.NoError(err)
	assert.NotEmpty(content)
}

func TestReplaceAllEntry(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)
	assert := assert.New(t)

	ReplaceAllEntry(&node_root, "mysql:10.5", "mariadb:10.5")

	content, err := yaml.Marshal(&node_root)
	assert.NoError(err)
	assert.NotEmpty(content)
}
