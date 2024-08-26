package cmd

import (
	"fmt"
	"log"
	"strings"

	entity "github.com/upsun/lib-sun/entity"
	utils "github.com/upsun/lib-sun/utility"
)

func RunHook(ws string, hook string) {
	log.Print("Run hooks on GIT repository...")

	// for _, hook := range hooks {
	fmt.Println(hook)
	payload := strings.Split(hook, " ")
	cmd := payload[0]
	payload = payload[1:]
	stdout, stderr, err := utils.CallExePath(cmd, "", ws, payload...)

	// Why twice ?
	if err != nil {
		log.Println("Hook is not good !")
		log.Printf("error  : %v\n", err)
		log.Printf("stdout : %v\n", string(stdout))
		log.Printf("stderr : %v\n", stderr)
	} else {
		fmt.Print(string(stdout))
	}
	// }
}

func RunRemoteHook(projectContext entity.ProjectGlobal, app string, posthook string) {
	log.Print("Run hooks on remote project...")

	// for _, posthook := range posthooks {
	payload := []string{
		"--environment=" + projectContext.DefaultEnv,
		"--app=" + app,
		posthook,
	}
	output, err := utils.CallCLIString(projectContext, "ssh", payload...)
	if err != nil {
		log.Printf("command execution failed: %v \n", err)
	} else {
		fmt.Print(output)
	}
	// }
}
