package cmd

import (
	"fmt"
	"log"
	"strings"

	app "github.com/upsun/lib-sun"
	entity "github.com/upsun/lib-sun/entity"
	utils "github.com/upsun/lib-sun/utility"
)

func sshPath(sshUrl string, path string) string {
	return sshUrl + ":/app" + path + "/"
}

func getSshUrl(projectCtx entity.ProjectGlobal, mount entity.EnvMount) string {
	payload := []string{
		"--environment=" + projectCtx.DefaultEnv,
		"--app=" + mount.Application,
		"--pipe",
	}
	sshUrl, err := utils.CallCLIString(projectCtx, "ssh", payload...)
	if err != nil {
		log.Printf("command execution failed: %v \n", err)
	}

	return sshUrl
}

func MountsExport(projectContext entity.ProjectGlobal, ws utils.PathTmp) {
	log.Print("Export Mounts...")
	_, env := utils.HackAlias(projectContext.Provider)

	for key, mount := range projectContext.Mounts {
		folder := ws.MakeDataFolder(key) + "/"
		log.Printf("Dump data for %v on %v", mount.Path, mount.Application)

		sshUrl := getSshUrl(projectContext, mount)
		payload := []string{
			"-azP",
			sshPath(sshUrl, mount.Path),
			folder,
		}

		fmt.Printf("%v %v\n", "rsync", strings.Join(payload, " "))

		output, stderr, err := utils.CallExe("rsync", env, payload...)
		if err != nil {
			log.Printf("command execution failed: %v \n %v", err, stderr)
		}

		// Debug Result
		if app.Args.Verbose {
			fmt.Printf("stdout : %v\n", string(output))
			fmt.Printf("stderr : %v\n%v\n", err, stderr)
		}

		mount.DumpPath = folder
		projectContext.Mounts[key] = mount
	}
}

func MountsImport(projectContext entity.ProjectGlobal) {
	log.Print("Import Mounts...")

	_, env := utils.HackAlias(projectContext.Provider)

	for _, mount := range projectContext.Mounts {
		log.Printf("Restore dump data for %v on %v", mount.Path, mount.Application)

		sshUrl := getSshUrl(projectContext, mount)
		payload := []string{
			"-azP",
			mount.DumpPath,
			sshPath(sshUrl, mount.Path),
		}

		fmt.Printf("%v %v\n", "rsync", strings.Join(payload, " "))

		output, stderr, err := utils.CallExe("rsync", env, payload...)
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
