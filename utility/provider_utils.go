package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"

	app "github.com/upsun/lib-sun"
	entity "github.com/upsun/lib-sun/entity"
)

var PROVIDER_SUPPORTED = []string{
	entity.UPS_PROVIDER,
	entity.PSH_PROVIDER,
	"ibexa",
	"shopware",
	"pimcore-cloud",
	"webpaas",
}

var PROVIDER_DOMAIN = map[string]string{
	entity.UPS_PROVIDER: "upsun.com",
	entity.PSH_PROVIDER: "platform.sh",
}

type PshResult struct {
	Status  any    `json:"status"`
	Code    int    `json:"code,omitempty"`
	Title   string `json:"title"`
	Message string `json:"message,omitempty"`
	// Detail  struct `json:"detail,omitempty"` {  // Not common to all result
	// 	Error string
	// }
}

func (c *PshResult) Clean() {
	c.Status = ""
	c.Code = 0
	c.Title = ""
	c.Message = ""
	//c.Detail.Error = ""  // Not common to all result
}

func CallExePath(name string, env string, path string, payload ...string) ([]byte, string, error) {
	var stderrBuffer bytes.Buffer

	cmd := exec.Command(name, payload...)
	cmd.Stderr = &stderrBuffer
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, env)
	if path != "" {
		cmd.Dir = path
	}

	stdout, err := cmd.Output()
	stderr := stderrBuffer.String()

	// Debug Result
	if app.Args.Verbose || err != nil {
		fmt.Printf("error  : %v\n", err)
		fmt.Printf("stdout : %v\n", string(stdout))
		fmt.Printf("stderr : %v\n", stderr)
	}

	return stdout, stderr, err
}

func CallExe(name string, env string, payload ...string) ([]byte, string, error) {
	return CallExePath(name, env, "", payload...)
}

func CallCLIPath(project entity.ProjectGlobal, command string, path string, payload ...string) ([]byte, error) {
	provider, env := HackAlias(project.Provider)

	if project.ID != "" {
		payload = append([]string{command, "-p", project.ID}, payload...)
	} else {
		payload = append([]string{command}, payload...)
	}

	if app.Args.Verbose {
		payload = append([]string{"-vvv"}, payload...)
	}
	fmt.Println(project.Provider + " " + strings.Join(payload, " "))
	output, _, err := CallExePath(provider, env, path, payload...)

	return output, err
}

func CallCLI(project entity.ProjectGlobal, command string, payload ...string) ([]byte, error) {
	return CallCLIPath(project, command, "", payload...)
}

func CallAPI(projectContext entity.ProjectGlobal, payload ...string) PshResult {
	var result PshResult

	output, err := CallCLI(projectContext, "project:curl", payload...)
	if err != nil {
		log.Printf("command execution failed: %s", err)
	}

	if len(output) > 0 {
		if err := json.Unmarshal(output, &result); err != nil {
			log.Printf("failed to unmarshal response: %s", err)
		}
		fmt.Printf("\tDone : %v [%v] %v\n", result.Status, result.Code, result.Message)
	}

	return result
}

func HackAlias(providerRequest string) (string, string) {
	// Hack golang not found bash alias.
	provider := providerRequest
	env := ""

	_, err := exec.LookPath(providerRequest)
	if err != nil {
		provider = entity.PSH_PROVIDER

		// TODO Check if binary exist before use hack from alias.
		home, _ := os.UserHomeDir()
		switch providerRequest {
		case entity.UPS_PROVIDER:
			env = "CLI_CONFIG_FILE=" + home + "/.platformsh/cli-config/upsun-cli.yaml"
		case "shopware":
			env = "CLI_CONFIG_FILE=" + home + "/.platformsh/cli-config/shopware-cli.yaml"
		case "ibexa":
			env = "CLI_CONFIG_FILE=" + home + "/.platformsh/cli-config/ibexa-cli.yaml"
		case "magento-cloud":
			env = "CLI_CONFIG_FILE=" + home + "/.platformsh/cli-config/magento-cloud-cli.yaml"
		case "pimcore-cloud":
			env = "CLI_CONFIG_FILE=" + home + "/.platformsh/cli-config/pimcore-cli.yaml"
		case "webpaas":
			env = "CLI_CONFIG_FILE=" + home + "/.platformsh/cli-config/webpaas-cli.yaml"
		}
	}
	// End Hack

	return provider, env
}

func CliJsonString(key string, value string) string {
	return "\"" + key + "\":\"" + value + "\""
}

func CliJsonBoolean(key string, value string) string {
	return "\"" + key + "\":" + value
}

func CliJson(values ...string) string {
	var sb strings.Builder
	var last = len(values) - 1

	sb.WriteString("{")
	for i, value := range values {
		sb.WriteString(value)
		if i < last {
			sb.WriteString(",")
		}
	}
	sb.WriteString("}")

	return sb.String()
}

func CallCLIString(project entity.ProjectGlobal, command string, payload ...string) (string, error) {
	output, err := CallCLI(project, command, payload...)
	if err != nil {
		return "", err
	}
	return string(output[:]), nil
}

// func Authenticated(project entity.ProjectGlobal) {
// 	payload := []string{"--force", "--yes"}

// 	projectEmpty := project
// 	projectEmpty.ID = ""
// 	output, err := CallCLI(projectEmpty, "auth:browser-login", payload...)

// 	fmt.Println(output, err)
// }

func IsAuthenticated(project entity.ProjectGlobal) bool {
	log.Printf("Check Authentification for %v...", project.Provider)
	payload := []string{"--refresh", "--format=csv", "--yes"}

	projectEmpty := project
	projectEmpty.ID = ""

	output, err := CallCLIString(projectEmpty, "auth:info", payload...)
	if err != nil {
		return false
	}

	var (
		id     string
		email  string
		result bool
	)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "id") {
			id = strings.Split(line, ",")[1]
		}
		if strings.Contains(line, "email") {
			email = strings.Split(line, ",")[1]
			result = true
		}
	}

	if id != "" {
		project.Users[id] = entity.ProjectUser{
			Id:    id,
			Email: email,
		}
	}

	return result
}

func HasSufficientRights(provider string) (bool, error) {
	cmd := exec.Command(provider, "project:info")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return strings.Contains(string(output), "Access"), nil
}

func ProviderCheck(provider string) {
	// Normalize the input to lowercase to handle case-insensitivity
	provider = strings.ToLower(provider)

	// Check for valid input variants
	if !slices.Contains(PROVIDER_SUPPORTED, provider) {
		fmt.Println("Invalid CLI provider entered. Please enter 'upsun' or 'platform'")
		os.Exit(-1)
	}
}
