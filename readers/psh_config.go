package readers

import (
	"fmt"
	"os"
	"strings"

	entity "github.com/upsun/lib-sun/entity"
	"gopkg.in/yaml.v3"
)

const (
	KEY_NAME = "name"
	KEY_SRC  = "source"
	KEY_ROOT = "root"
)

// Internal function to load Yaml file on MetaModel.
//
//   - node: Metaconfig node to update.
//   - filePath: Path of Yaml file.
func readConfigFile(node *yaml.Node, filePath string) {
	// Read File Services.yaml .
	data_service, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		os.Exit(1)
	}

	var sub_metamodel yaml.Node
	if len(data_service) > 0 {
		// Unmarshal yaml file to node graph.

		err = yaml.Unmarshal(data_service, &sub_metamodel)
		if err != nil {
			fmt.Printf("Error parsing YAML file: %s\n", err)
			os.Exit(1)
		}

		if sub_metamodel.Kind > 0 && len(sub_metamodel.Content) > 0 {
			// Config node type.
			node.Kind = yaml.MappingNode
			node.Tag = TAG_MAP

			// Map to meta-model.
			node.Content = append(node.Content, sub_metamodel.Content[0].Content...)
		} else {
			fmt.Printf("No content to append from %s\n", filePath)
		}

	} else {
		fmt.Printf("No content to append from %s\n", filePath)
		// Map to meta-model.
		//node.Content = append(node.Content, "")
	}
}

func ReadServices(metamodel *entity.MetaConfig, filePaths []string) {
	for _, filePath := range filePaths {
		readConfigFile(&metamodel.Services, filePath)
	}
}

func ReadApplications(metamodel *entity.MetaConfig, filePaths []string, projectWorkspace string) {
	for _, filePath := range filePaths {
		var read_node yaml.Node

		readConfigFile(&read_node, filePath)
		for _, appLevel := range read_node.Content {
			var header_node yaml.Node
			nameEntry := FindEntry(appLevel, KEY_NAME)

			//  Name
			header_node.Kind = yaml.ScalarNode
			header_node.Tag = TAG_STRING
			header_node.Value = nameEntry.ValueNode.Value

			//  Value
			value_node := appLevel
			RemoveEntry(value_node, KEY_NAME)

			// Merge all
			metamodel.Applications.Kind = yaml.MappingNode
			metamodel.Applications.Tag = TAG_MAP
			metamodel.Applications.Content = append(metamodel.Applications.Content, &header_node, value_node)
		}
	}
}

func ReadPlatforms(metamodel *entity.MetaConfig, filePaths []string, projectWorkspace string) {
	for _, filePath := range filePaths {
		var (
			read_node   yaml.Node
			header_node yaml.Node
		)

		readConfigFile(&read_node, filePath)
		entry := FindEntry(&read_node, KEY_NAME)

		// Extract App section
		//  Name
		header_node.Kind = yaml.ScalarNode
		header_node.Tag = TAG_STRING
		header_node.Value = entry.ValueNode.Value

		//  Value
		value_node := &read_node
		RemoveEntry(value_node, KEY_NAME)

		// Add project path
		rootProject := strings.ReplaceAll(filePath, projectWorkspace, "")
		rootProject = strings.ReplaceAll(rootProject, entity.PSH_PLATFORM, "")

		rootpath := CreateEntryString(KEY_ROOT, rootProject)

		source := FindEntryFirstLevel(&read_node, KEY_SRC)
		if source.KeyNode == nil {
			// Create Source node and add Root
			source = CreateEntryEntity(KEY_SRC, rootpath)
			AddEntry(value_node, source)
		} else {
			// Append Root to Source node
			rootpathFinded := FindEntryFirstLevel(source.ValueNode, KEY_ROOT)
			if rootpathFinded.KeyNode == nil {
				AddEntry(source.ValueNode, rootpath)
			}
		}

		// Merge all
		metamodel.Applications.Kind = yaml.MappingNode
		metamodel.Applications.Tag = TAG_MAP
		metamodel.Applications.Content = append(metamodel.Applications.Content, &header_node, value_node)
	}
}

func ReadRoutes(metamodel *entity.MetaConfig, filePaths []string) {
	for _, filePath := range filePaths {
		readConfigFile(&metamodel.Routes, filePath)
	}
}
