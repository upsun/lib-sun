package cmd

import (
	"fmt"
	"log"
	"strings"

	utils "github.com/upsun/lib-sun/utility"
)

func RunHook(ws string, hooks ...string) {
	log.Print("Run hooks on GIT repository...")
	for _, hook := range hooks {
		fmt.Println(hook)
		payload := strings.Split(hook, " ")
		cmd := payload[0]
		payload = payload[1:]
		stdout, stderr, err := utils.CallExePath(cmd, "", ws, payload...)
		if err != nil {
			log.Println("Hook is not good !")
			log.Printf("error  : %v\n", err)
			log.Printf("stdout : %v\n", string(stdout))
			log.Printf("stderr : %v\n", stderr)
		}
	}
}
