package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
	app "github.com/upsun/lib-sun"
	entity "github.com/upsun/lib-sun/entity"
)

func Disclaimer(appName string) {
	if !app.Args.Silent {
		fmt.Print(appName + " from ADV-initial version " + app.VERSION + "\n\n")
	}
}

func LinkToProject(projectContext entity.ProjectGlobal) {
	if !app.Args.Silent {
		link := fmt.Sprintf("https://console.%s/%v/%v", PROVIDER_DOMAIN[projectContext.Provider], projectContext.OrgEmail, projectContext.ID)
		log.Printf("Go to project: %v", link)
	}
}

func IsFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func RequireFlag(name string, question string, defaultValue string, block bool) string {
	result := defaultValue

	if !IsFlagPassed(name) {
		for {
			fmt.Printf(question, defaultValue)

			_, err := fmt.Scanln(&result)
			if err != nil && !strings.Contains(err.Error(), "unexpected newline") {
				fmt.Printf("Error : %q", err)
				os.Exit(-1)
			}

			// Use default value
			if result == "" {
				result = defaultValue
			}

			if result != "" {
				block = false
			}

			if !block {
				break
			}
		}
	}
	return result
}

func InitLogger(appName string) {
	var mw io.Writer

	// defer func() {
	// 	cerr := logFile.Close()
	// 	if err == nil {
	// 		err = cerr
	// 	}
	// }()

	if app.Args.Silent {
		// Hack need case without log file
		mw = io.MultiWriter(os.Stdout)
	} else {
		logFile, err := os.OpenFile(fmt.Sprintf(app.Args.PathLog+"log-%v-%v.log", appName, TimeStamp()), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
		mw = io.MultiWriter(os.Stdout, logFile)
	}

	log.SetOutput(mw)
	log.Println(strings.Join(os.Args, " "))
	log.Println()
}
