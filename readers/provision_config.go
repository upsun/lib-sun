package readers

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
	entity "upsun.com/lib-upsun/entity"
	utils "upsun.com/lib-upsun/utility"
)

func ReadProvisionFile(projectDstContext *entity.ProjectGlobal, provisionContext *entity.ProvisionGlobal, configPath string, ws utils.PathTmp) {
	log.Print("Read Provision metadata...")

	// Collect Meta-model & Data
	filePath := path.Join(configPath, "metadata.yaml")
	data_service, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		os.Exit(1)
	}
	log.Printf("Provision metadata found at: %v", filePath)

	// Unmarshal yaml file to node graph.
	err = yaml.Unmarshal(data_service, &provisionContext)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
		os.Exit(1)
	}

	// Populate ProjectContext
	projectDstContext.Name = provisionContext.Name
	projectDstContext.Description = provisionContext.Description
	projectDstContext.Timezone = "Etc/UTC"

	// Variables
	for k, v := range provisionContext.Variables {
		log.Printf("Find variable: %q", k)

		variable := entity.ProjectVariable{
			Name:        k,
			Value:       v,
			IsAtBuild:   true,
			IsAtDeploy:  true,
			IsSensitive: false,
		}
		projectDstContext.Variables[variable.Name] = variable
	}
	// Access
	for k, v := range provisionContext.Users {
		log.Printf("Find Access: %q", k)

		access := entity.ProjectAccess{
			UserId:      k,
			Permissions: strings.Split(v, ","),
		}
		projectDstContext.Access[access.UserId] = access
	}

	for app_name, app := range provisionContext.Applications {
		// Update struct
		app.Name = app_name
		provisionContext.Applications[app_name] = app

		// Load Services
		for rel, src := range app.Services {
			log.Printf("Find Service: %q [Type: %v]", rel, "Not needed")

			absDst := ws.MakeDataFolder(rel) + "/"
			filePath := path.Join(absDst, utils.TimeStamp()+".sql")
			utils.GetFile(src, filePath, configPath)

			service := entity.EnvService{
				Type:         "mariadb", // Not need to know for restore !
				DumpPath:     filePath,
				Application:  app.Name,
				Relationship: rel,
			}
			projectDstContext.Services[rel] = service
		}

		// Load Mounts
		for k, v := range app.Mounts {
			log.Printf("Find Mount : %v => %v", k, v)

			dst := ws.MakeDataFolder(k) + "/"
			utils.GetFile(v, dst, configPath)

			mount := entity.EnvMount{
				DumpPath:    dst,
				Path:        k,
				Application: app.Name,
			}
			projectDstContext.Mounts[k] = mount
		}
	}

}
