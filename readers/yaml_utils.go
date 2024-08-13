package readers

import (
	"log"

	"gopkg.in/yaml.v3"
)

const (
	TAG_STRING = "!!str"
	TAG_MAP    = "!!map"
	MAX_RECURS = 30
)

type KeyValueNodes struct {
	KeyNode   *yaml.Node
	ValueNode *yaml.Node
}

func NewKeyValueNodes() KeyValueNodes {
	result := KeyValueNodes{}
	result.KeyNode = &yaml.Node{}
	result.ValueNode = &yaml.Node{}

	return result
}

func CreateEntryEmpty(keyName string) KeyValueNodes {
	result := NewKeyValueNodes()

	// Key
	result.KeyNode.Kind = yaml.ScalarNode
	result.KeyNode.Tag = TAG_STRING
	result.KeyNode.Value = keyName

	return result
}

func CreateEntryString(keyName string, valueContent string) KeyValueNodes {
	result := CreateEntryEmpty(keyName)

	// Value
	result.ValueNode.Kind = yaml.ScalarNode
	result.ValueNode.Tag = TAG_STRING
	result.ValueNode.Value = valueContent

	return result
}

func CreateEntryEntity(keyName string, valueContent KeyValueNodes) KeyValueNodes {
	result := CreateEntryEmpty(keyName)

	// Value
	result.ValueNode.Kind = yaml.MappingNode
	result.ValueNode.Tag = TAG_MAP
	AddEntry(result.ValueNode, valueContent)

	return result
}

func AddEntry(node *yaml.Node, entry KeyValueNodes) {
	node.Content = append(node.Content, entry.KeyNode, entry.ValueNode)
}

func FindEntry(node *yaml.Node, keyFinded string) KeyValueNodes {
	nodeKey, nodeValue, _, _ := findChildNode(node, keyFinded, 0, MAX_RECURS)

	result := NewKeyValueNodes()
	result.KeyNode = nodeKey
	result.ValueNode = nodeValue

	return result
}

func FindEntryFirstLevel(node *yaml.Node, keyFinded string) KeyValueNodes {
	nodeKey, nodeValue, _, _ := findChildNode(node, keyFinded, 0, 1)

	result := NewKeyValueNodes()
	result.KeyNode = nodeKey
	result.ValueNode = nodeValue

	return result
}

// Find a node by key.
// Return key, value, parent, index.
func findChildNode(node *yaml.Node, keyFinded string, level int, levelLimit int) (*yaml.Node, *yaml.Node, *yaml.Node, int) {
	for i := 0; i < len(node.Content); i++ {
		nodeKey := node.Content[i]

		if nodeKey.Value == keyFinded {
			return nodeKey, node.Content[i+1], node, i
		}

		if level < levelLimit {
			if nodeKey, nodeValue, nodeParent, index := findChildNode(nodeKey, keyFinded, level+1, levelLimit); nodeKey != nil {
				return nodeKey, nodeValue, nodeParent, index
			}
		}
	}
	return nil, nil, nil, -1
}

func RemoveEntry(node *yaml.Node, keyFinded string) *yaml.Node {
	_, _, nodeParent, index := findChildNode(node, keyFinded, 0, MAX_RECURS)

	if nodeParent != nil {
		content := nodeParent.Content
		nodeParent.Content = append(content[:index], content[index+2:]...)
		return nodeParent
	} else {
		return nil
	}
}

func RemoveAllEntry(node *yaml.Node, keyFinded string) {
	owner := RemoveEntry(node, keyFinded)

	for owner != nil {
		log.Printf("- Remove '%s' field on line %v", keyFinded, owner.Line) // TODO(mick): return node key of parent
		owner = RemoveEntry(node, keyFinded)
	}
}

func ReplaceEntry(node *yaml.Node, find string, replace string) *yaml.Node {
	nodeFinded, _, _, _ := findChildNode(node, find, 0, MAX_RECURS)

	if nodeFinded != nil {
		nodeFinded.Value = replace
		return nodeFinded
	}

	return nil
}

func ReplaceAllEntry(node *yaml.Node, find string, replace string) {
	owner := ReplaceEntry(node, find, replace)

	for owner != nil {
		log.Printf("- Replace '%v' by '%v' line %v", find, replace, owner.Line) // TODO(mick): return node key of parent
		owner = ReplaceEntry(node, find, replace)
	}
}
