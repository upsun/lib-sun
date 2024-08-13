package writers

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
	entity "upsun.com/lib-upsun/entity"
)

func GenerateUpsunConfigFile(metamodel entity.MetaConfig, outputFilePath string) {
	file_config, err := yaml.Marshal(&metamodel)
	if err != nil {
		fmt.Printf("Error marshalling YAML: %s\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(outputFilePath, file_config, 0644)
	if err != nil {
		fmt.Printf("Error writing YAML file: %s\n", err)
		os.Exit(1)
	}
}
