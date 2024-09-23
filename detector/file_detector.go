package detector

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	entity "github.com/upsun/lib-sun/entity"
)

const PSH_CONFIG_PATH = ".platform/"

// searchFile recursively searches for a file with the specified name in the given directory and its subdirectories.
func searchFile(dir, targetFile string) ([]string, error) {
	var filePaths []string

	// Walk through the directory and its subdirectories
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// Check if the current path is a regular file and matches the target file name
		if err == nil && !info.IsDir() && info.Name() == targetFile {
			filePaths = append(filePaths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(filePaths) == 0 {
		return nil, fmt.Errorf("file '%s' not found in directory '%s'", targetFile, dir)
	}

	return filePaths, nil
}

func FindConfig(searchDir string) (map[string][]string, error) {
	// Directory to search in (change this to the directory you want to search in)
	foundFiles := map[string][]string{}

	// Validate if the entered directory exists
	if _, err := os.Stat(searchDir); os.IsNotExist(err) {
		log.Printf("Directory '%s' does not exist\n", searchDir)
		return nil, err
	}

	// Target files to search for
	filesToCheck := map[string]string{
		entity.PSH_ROUTE:       PSH_CONFIG_PATH,
		entity.PSH_SERVICE:     PSH_CONFIG_PATH,
		entity.PSH_APPLICATION: PSH_CONFIG_PATH,
		entity.PSH_PLATFORM:    "",
	}

	// Search for the target files in the directory and its subdirectories
	for file, dir := range filesToCheck {
		filePaths, err := searchFile(filepath.Join(searchDir, dir), file)
		if err != nil {
			fmt.Println("WARNING:", err)
			continue
		}
		for _, filePath := range filePaths {
			log.Printf("Discover found '%s' at: %s\n", file, filePath)
		}
		foundFiles[file] = filePaths
	}

	return foundFiles, nil
}
