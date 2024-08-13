package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"path"
	"slices"
	"strings"

	app "github.com/upsun/lib-upsun"
	entity "github.com/upsun/lib-upsun/entity"
	utils "github.com/upsun/lib-upsun/utility"
)

var ServicesDumpSupported = []string{
	"mysql",
	"mariadb",
	"oracle-mysql",
	"postgresql",
}

type relationship struct {
	Service string
}

type webapp struct {
	Relationships map[string]relationship
	Mounts        map[string]entity.EnvMount
}

type envResult struct {
	Services map[string]entity.EnvService
	Webapps  map[string]webapp
}

func ServicesMountsRead(projectContext entity.ProjectGlobal) {
	log.Print("Read Services & Mounts...")

	payload := []string{"-X", "GET", "/environments/" + projectContext.DefaultEnv + "/deployments/current"}
	jsonContent, err := utils.CallCLI(projectContext, "project:curl", payload...)
	if err != nil {
		log.Printf("command execution failed: %s", err)
	}

	var variables envResult
	if err := json.Unmarshal(jsonContent, &variables); err != nil {
		log.Printf("failed to unmarshal response: %s", err)
	}

	// Services Read
	for nameSrv, service := range variables.Services {
		log.Printf("Find Service: %q [Type: %v]", nameSrv, service.Type)

		for nameApp, application := range variables.Webapps {
			for nameRel, rel := range application.Relationships {
				if nameSrv == rel.Service {
					srv := variables.Services[nameSrv]
					srv.Relationship = nameRel
					srv.Application = nameApp
					srv.TypeService = rel.Service
					variables.Services[nameSrv] = srv
					break
				}
			}
		}
	}
	maps.Copy(projectContext.Services, variables.Services)

	// Mounts Read
	for nameApp, application := range variables.Webapps {
		for nameMnt, mnt := range application.Mounts {
			key := nameApp + "/" + nameMnt
			// TODO Check if mount exist
			//if _, ok := projectContext.Mounts[key]; !ok {
			log.Printf("Find Mount : %v => %v.%v [%v]", nameMnt, nameApp, mnt.SourcePath, mnt.Type)
			mnt.Application = nameApp
			mnt.Path = nameMnt
			projectContext.Mounts[key] = mnt
			//}
		}
	}
}

func DisplayUnsupportedServices(projectContext entity.ProjectGlobal) {

	log.Printf("upsun_clone do not clone some services ! Please dump/restore them manualy for :")
	for nameSer, service := range projectContext.Services {
		typeService := strings.Split(service.Type, ":")[0]

		if !slices.Contains(ServicesDumpSupported, typeService) {
			log.Printf("\t%v [%v]", nameSer, service.Type)
		}
	}
}

func ServicesExport(projectContext entity.ProjectGlobal, ws utils.PathTmp) {
	log.Print("Export Services...")

	for name, service := range projectContext.Services {
		typeService := strings.Split(service.Type, ":")[0]

		if slices.Contains(ServicesDumpSupported, typeService) {
			folder := ws.MakeDataFolder(name)
			filePath := path.Join(folder, utils.TimeStamp()+".sql.gz")

			log.Printf("Dump data for %v on %v", name, filePath)

			payload := []string{
				"--file=" + filePath,
				"--gzip",
				"--environment=" + projectContext.DefaultEnv,
				"--app=" + service.Application,
				"--relationship=" + service.Relationship,
				"--yes",
				"--no-interaction",
			}
			_, err := utils.CallCLIString(projectContext, "db:dump", payload...)
			if err != nil {
				log.Printf("command execution failed: %s", err)
			}

			service.DumpPath = filePath
			projectContext.Services[name] = service

			if utils.IsExist(filePath) {
				log.Print("Database dump downloaded.\n")
			} else {
				log.Print("The database dump wasn't created.\n")
			}
		}
	}
}

func ServicesImport(projectContext entity.ProjectGlobal) {
	log.Print("Import Services...")
	var payload []string
	provider, env := utils.HackAlias(projectContext.Provider)

	for name, service := range projectContext.Services {
		typeService := strings.Split(service.Type, ":")[0]

		if slices.Contains(ServicesDumpSupported, typeService) {
			log.Printf("Import dump data for %v on %v", name, service.DumpPath)

			if strings.HasSuffix(service.DumpPath, ".sql.gz") {
				payload = []string{
					"-c",
					"gunzip < " + service.DumpPath +
						" | " + provider + " -vvv " +
						" db:sql" +
						" --project=" + projectContext.ID +
						" --environment=" + projectContext.DefaultEnv +
						" --app=" + service.Application +
						" --relationship=" + service.Relationship +
						" --yes" +
						" --no-interaction",
				}
			} else {
				payload = []string{
					"-c",
					"cat " + service.DumpPath +
						" | " + provider + " -vvv " +
						" db:sql" +
						" --project=" + projectContext.ID +
						" --environment=" + projectContext.DefaultEnv +
						" --app=" + service.Application +
						" --relationship=" + service.Relationship +
						" --yes" +
						" --no-interaction",
				}
			}

			// if version.DEBUG {
			fmt.Println(payload[1])
			// }

			output, stderr, err := utils.CallExe("bash", env, payload...)
			if err != nil {
				log.Printf("command execution failed: %v \n %v", err, stderr)
			}

			// Debug Result
			if app.Args.Verbose {
				fmt.Printf("stdout : %v\n", string(output))
				fmt.Printf("stderr : %v\n%v\n", err, stderr)
			}

		}
	}
}
