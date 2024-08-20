package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	app "github.com/upsun/lib-sun"
	entity "github.com/upsun/lib-sun/entity"
	utils "github.com/upsun/lib-sun/utility"
)

const REPO_BIN = "git"

func ImportRepository(projectContext entity.ProjectGlobal, wsFolder string) {
	log.Print("Pull Repository...")

	if app.ArgsM.PshRepo {
		url := fmt.Sprintf("%v@git.%v:%v.git", projectContext.ID, projectContext.Region, projectContext.ID)
		payload := append([]string{"clone", "--recursive", "--branch"}, projectContext.DefaultEnv, url, wsFolder)

		callVCS(wsFolder, payload...)
	} else {
		os.RemoveAll(wsFolder)
		payload := []string{
			"--environment=" + projectContext.DefaultEnv,
			"--no-interaction",
			"--yes",
			projectContext.ID,
			wsFolder,
		}

		copy := projectContext // Clone
		copy.ID = ""           // Cli call without -p argument.
		_, err := utils.CallCLI(copy, "project:get", payload...)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Remove remote repository
	payload := append([]string{"remote", "remove"}, projectContext.Provider)
	callVCS(wsFolder, payload...)
}

func ShallowCloneRepository(url string, version string, wsFolder string) {
	log.Print("Clone Repository...")

	err := os.MkdirAll(wsFolder, 0750)
	if err != nil {
		log.Fatal(err)
	}

	payload := append([]string{"clone", "--recursive", "--depth", "1", "--single-branch", "--branch"}, version, url, wsFolder)
	callVCS(wsFolder, payload...)

	// Re-init Git shallow repository.
	os.RemoveAll(path.Join(wsFolder, ".git"))
}

func InitRepository(projectContext entity.ProjectGlobal, wsFolder string) {
	log.Print("Initialize GIT repository...")
	payload := []string{
		"init",
		"--initial-branch=" + projectContext.DefaultEnv,
	}
	callVCS(wsFolder, payload...)
}

func ExportRepository(projectContext entity.ProjectGlobal, wsFolder string) {
	log.Print("Push Repository... (please wait, it can be long ! see on activity)")

	copy := projectContext // Clone
	copy.ID = ""           // Cli call without -p argument.

	payload := []string{
		//"--environment=" + projectContext.DefaultEnv,
		"--no-interaction",
		"--yes",
		projectContext.ID,
	}

	_, err := utils.CallCLIPath(copy, "project:set-remote", wsFolder, payload...)
	if err != nil {
		log.Fatal(err)
	}

	payload = []string{
		//"--target=" + projectContext.DefaultEnv,
		"--no-interaction",
		"--yes",
	}

	_, err = utils.CallCLIPath(copy, "environment:push", wsFolder, payload...)
	if err != nil {
		log.Fatal(err)
	}
}

func CommitOnRepository(wsFolder string, files ...string) {
	log.Print("Commit files on GIT repository...")

	// Add to git index.
	payload := append([]string{"add"}, files...)
	callVCS(wsFolder, payload...)

	// Commit
	payload = []string{
		"commit",
		"-m",
		"Add PaaS configuration files.",
	}
	callVCS(wsFolder, payload...)
}

func CommitAllOnRepository(wsFolder string, msg string) {
	log.Print("Commit all files on GIT repository...")

	// Add to git index.
	payload := []string{"add", wsFolder}
	callVCS(wsFolder, payload...)

	// Commit
	payload = []string{
		"commit",
		"-m",
		msg,
	}
	callVCS(wsFolder, payload...)
}

func AddToRepository(files map[string]string, wsFolder string, configPath string) {
	log.Print("Add files to repository...")

	for src, dst := range files {
		absDstPath := path.Join(wsFolder, dst)
		dstFolder := filepath.Dir(absDstPath)

		log.Printf("Add file %v in %v", src, absDstPath)
		err := os.MkdirAll(dstFolder, 0750)
		if err != nil {
			fmt.Println(err)
		}

		utils.GetFile(src, absDstPath, configPath)
	}
}

func AddConfigToRepository(srcFolder string, wsFolder string, app string) {
	log.Print("Add Paas config on repository...")

	pshSrcConfigPath := path.Join(srcFolder, "psh", app)

	if utils.IsExist(pshSrcConfigPath) {
		pshDstAppPath, pshDstPath := utils.MakePshConfigPath(wsFolder, app)
		err := utils.CopyFile(path.Join(pshSrcConfigPath, entity.PSH_PLATFORM), pshDstAppPath)
		if err != nil {
			log.Fatal(err)
		}

		err = os.MkdirAll(pshDstPath, 0750)
		if err != nil {
			log.Fatal(err)
		}
		err = utils.CopyFile(path.Join(pshSrcConfigPath, entity.PSH_SERVICE), path.Join(pshDstPath, entity.PSH_SERVICE))
		if err != nil {
			log.Fatal(err)
		}
		err = utils.CopyFile(path.Join(pshSrcConfigPath, entity.PSH_ROUTE), path.Join(pshDstPath, entity.PSH_ROUTE))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("No configuration found ! (can be generate by 3th party ?)")
	}
}

func callVCS(wsFolder string, payload ...string) {
	fmt.Println(REPO_BIN + " " + strings.Join(payload, " "))

	_, _, err := utils.CallExePath(REPO_BIN, "", wsFolder, payload...)
	if err != nil {
		log.Fatal(err)
	}
}
