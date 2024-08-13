package cmd

import (
	"encoding/json"
	"log"
	"strings"

	entity "github.com/upsun/lib-upsun/entity"
	utils "github.com/upsun/lib-upsun/utility"
)

func ProjectCreate(projectContext *entity.ProjectGlobal) {
	log.Print("Create Project... (please wait, it can be long !)")

	// Security check
	if projectContext.ID != "" {
		log.Println("Project with ID cannot be created !")
		return
	}

	payload := []string{
		"--org=" + projectContext.OrgEmail,
		"--title=" + projectContext.Name,
		"--region=" + projectContext.Region,
		"--default-branch=" + projectContext.DefaultEnv,
		"--environments=" + "3", // TODO TBD
		"--storage=" + "5",      // TODO TBD
		"--no-set-remote",
		"--no-interaction",
		"--yes",
	}
	output, err := utils.CallCLIString(*projectContext, "project:create", payload...)
	if err != nil {
		log.Printf("command execution failed: %s", err)
	}

	// Bad code !
	projectContext.ID = strings.ReplaceAll(output, "\n", "")
}

func ProjectRead(projectContext *entity.ProjectGlobal) error {
	log.Print("Read Project...")

	payload := []string{"-X", "GET", "/"}
	jsonContent, err := utils.CallCLI(*projectContext, "project:curl", payload...)
	if err != nil {
		log.Printf("command execution failed: %s", err)
		return err
	}

	if err := json.Unmarshal(jsonContent, projectContext); err != nil {
		log.Printf("failed to unmarshal response: %s", err)
		return err
	}

	return nil
}

func ProjectWrite(projectContext entity.ProjectGlobal) {
	log.Print("Write Project...")

	payload := []string{"-X", "PATCH", "/", "--json",
		// DTO (Manual)
		utils.CliJson(
			utils.CliJsonString("title", projectContext.Name),
			utils.CliJsonString("description", projectContext.Description),
			utils.CliJsonString("timezone", projectContext.Timezone),
		),
	}

	// UPDATE project setting (base)
	utils.CallAPI(projectContext, payload...)
}
