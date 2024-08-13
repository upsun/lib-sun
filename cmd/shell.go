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
		_, _, err := utils.CallExePath(cmd, "", ws, payload...)
		if err != nil {
			log.Println("Hook is not good !")
		}
	}
}
