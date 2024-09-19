package cmd

import (
	"encoding/json"
	"log"
	"strings"

	entity "github.com/upsun/lib-sun/entity"
	utils "github.com/upsun/lib-sun/utility"
)

const LEVEL_PROJECT = "project"
const LEVEL_ENV = "environment"

func DisplaySensitiveVariables(projectContext entity.ProjectGlobal) {
	log.Printf("upsun_clone do not copy sensitive variables ! Please edit them manualy for :")
	log.Printf("- At project level :")
	for nameVar, variable := range projectContext.Variables {
		if variable.IsSensitive {
			log.Printf("\t%v", nameVar)
		}
	}
	log.Printf("- At environment level :")
	for nameVar, variable := range projectContext.VariablesEnv {
		if variable.IsSensitive {
			log.Printf("\t%v", nameVar)
		}
	}
}

func VariablesRead(projectContext entity.ProjectGlobal) {
	log.Print("Read variables (project Level)...")

	payload := []string{"-X", "GET", "/variables"}
	jsonContent, err := utils.CallCLI(projectContext, "project:curl", payload...)
	if err != nil {
		log.Printf("command execution failed: %s", err)
	}

	var variables []entity.ProjectVariable
	if err := json.Unmarshal(jsonContent, &variables); err != nil {
		log.Printf("failed to unmarshal response: %s", err)
	}

	// TODO : resync array (remove not use)
	for _, variable := range variables {
		log.Printf("Find variable: %q", variable.Name)
		projectContext.Variables[variable.ID] = variable
	}
}

func VariablesWrite(projectContext entity.ProjectGlobal) {
	log.Print("Write variables (project Level)...")

	var result utils.PshResult
	payloadBase := []string{"-X", "POST", "/variables", "--json"}

	for _, variable := range projectContext.Variables {
		log.Printf("Write variable: %q", variable.Name)

		// DTO (dynamic)
		dto := variable                // Make a copy (for not invalidate original)
		dto.ID = ""                    // Remove field on DTO
		dto.IsInheritable = nil        // Remove field on DTO
		output, _ := json.Marshal(dto) // Convert to JSON (by Marshal)
		dtoJson := string(output)      // Transform to string

		// CREATE case.
		payloadInsert := append(payloadBase, dtoJson)
		result = utils.CallAPI(projectContext, payloadInsert...)

		// UPDATE case.
		if result.Code == 409 && !variable.IsSensitive {
			payloadUpdate := []string{"-X", "PATCH",
				"/variables/" + variable.ID,
				"--json", dtoJson}
			utils.CallAPI(projectContext, payloadUpdate...)
		}
	}

	// TODO : Remove old variables
}

func VariablesEnvRead(projectContext entity.ProjectGlobal) {
	log.Print("Read variables (environment Level)...")

	payload := []string{"-X", "GET", "/environments/" + projectContext.DefaultEnv + "/variables"}
	jsonContent, err := utils.CallCLI(projectContext, "project:curl", payload...)
	if err != nil {
		log.Printf("command execution failed: %s", err)
	}

	var variables []entity.ProjectVariable
	if err := json.Unmarshal(jsonContent, &variables); err != nil {
		log.Printf("failed to unmarshal response: %s", err)
	}

	// TODO : resync array (remove not use)
	for _, variable := range variables {
		log.Printf("Find variable: %q", variable.Name)
		projectContext.VariablesEnv[variable.ID] = variable
	}
}

func VariablesEnvWrite(projectContext entity.ProjectGlobal) {
	log.Print("Write variables (environment Level)...")

	var result utils.PshResult
	payloadBase := []string{"-X", "POST", "/environments/" + projectContext.DefaultEnv + "/variables", "--json"}

	for _, variable := range projectContext.VariablesEnv {
		log.Printf("Write variable: %q", variable.Name)

		// DTO (dynamic)
		dto := variable                // Make a copy (for not invalidate original)
		dto.ID = ""                    // Remove field on DTO
		dto.IsInheritable = nil        // Remove field on DTO
		output, _ := json.Marshal(dto) // Convert to JSON (by Marshal)
		dtoJson := string(output)      // Transform to string

		// CREATE case.
		payloadInsert := append(payloadBase, dtoJson)
		result = utils.CallAPI(projectContext, payloadInsert...)

		// UPDATE case.
		if result.Code == 409 && !variable.IsSensitive {
			payloadUpdate := []string{"-X", "PATCH",
				"/variables/" + variable.ID,
				"--json", dtoJson}
			utils.CallAPI(projectContext, payloadUpdate...)
		}
	}

	// TODO : Remove old variables
}

func PopulateSensitive(projectContext entity.ProjectGlobal, appName string) {
	log.Print("Get sensitive value from SSH container...")

	payload := []string{
		"--environment=" + projectContext.DefaultEnv,
	}
	if appName != "" {
		payload = append(payload, "--app="+appName)
	}
	payload = append(payload, "export")

	output, err := utils.CallCLIString(projectContext, "ssh", payload...)
	if err != nil {
		log.Printf("command execution failed: %v \n", err)
	}

	envVars := make(map[string]string)
	envVars_str := strings.Split(output, "\n")
	for _, envVar_str := range envVars_str[:len(envVars_str)-1] {
		// Sanitize
		envVar_str_san := strings.Replace(envVar_str, "declare -x ", "", 1) // Remove 'declare -x ' from export command

		envVar := strings.SplitN(envVar_str_san, "=", 2)
		if len(envVar) == 2 {
			key := envVar[0]
			value := envVar[1]
			// Sanitize
			if strings.HasPrefix(value, "\"") { // Remove escape double quote
				value = value[1:]            // Remove start
				value = value[:len(value)-1] // Remove End
			} else if strings.HasPrefix(value, "'") { // Remove escape simple quote
				value = value[1:]            // Remove start
				value = value[:len(value)-1] // Remove End
			} else if strings.HasPrefix(value, "$'") { // Remove extquote case
				value = value[2:]            // Remove start
				value = value[:len(value)-1] // Remove End
			} else if strings.HasPrefix(value, "$\"") { // Remove extquote case
				value = value[3:]            // Remove start
				value = value[:len(value)-2] // Remove End
			}
			envVars[key] = value
		} else {
			envVars[envVar[0]] = ""
		}
	}

	sensitiveInjectValue(projectContext.Variables, envVars)
	sensitiveInjectValue(projectContext.VariablesEnv, envVars)
}

func sensitiveInjectValue(variables map[string]entity.ProjectVariable, sensitiveVariables map[string]string) {
	for _, variable := range variables {
		if variable.IsSensitive {
			shortName := strings.Replace(variable.Name, "env:", "", 1)
			value, found := sensitiveVariables[shortName]
			if found {
				variable.Value = value
				variables[variable.Name] = variable
			}
		}
	}
}
