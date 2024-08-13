package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	app "github.com/upsun/lib-sun"
	entity "github.com/upsun/lib-sun/entity"
)

const PATH_PREFIX = "cli-*-tmp"
const PATH_REPO = "repo"
const PATH_DATA = "data"
const PSH_CONFIG_PATH = ".platform"

type PathTmp struct {
	Root   string
	Repo   string
	Data   string
	Config string
}

///// PSH SECTION /////

func BuildTemporyWorkspace() PathTmp {
	log.Print("Make tempory workspace...")

	var err error
	var result PathTmp

	result.Root, err = os.MkdirTemp("", PATH_PREFIX)
	if err != nil {
		log.Fatal(err)
	}

	result.Repo = path.Join(result.Root, PATH_REPO)
	err = os.MkdirAll(result.Repo, 0750)
	if err != nil {
		log.Fatal(err)
	}
	result.Data = path.Join(result.Root, PATH_DATA)
	err = os.MkdirAll(result.Data, 0750)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Tempory workspace created at: %v", result.Root)

	return result
}

func BuildPersistWorkspace() PathTmp {
	log.Print("Make persitant workspace...")

	var err error
	var result PathTmp

	result.Root = "./data/"
	err = os.Mkdir(result.Root, 0750)
	if err != nil {
		log.Fatal(err)
	}

	result.Repo = path.Join(result.Root, PATH_REPO)
	result.Data = path.Join(result.Root, PATH_DATA)

	return result
}

func (i *PathTmp) MakeDataFolder(service string) string {
	dataFolderPath := path.Join(i.Data, service)
	log.Printf("Generate Data storage : %v => %v", service, dataFolderPath)

	err := os.MkdirAll(dataFolderPath, 0750)
	if err != nil {
		log.Fatal(err)
	}

	return dataFolderPath
}

func MakeUpsunConfigPath(rootPath string) (string, string) {
	log.Print("Generate Upsun structure...")

	upsunFolderPath := path.Join(rootPath, ".upsun")
	err := os.MkdirAll(upsunFolderPath, 0750)
	if err != nil {
		log.Fatal(err)
	}

	return upsunFolderPath, path.Join(upsunFolderPath, "config.yaml")
}

func MakePshConfigPath(rootPath string, app string) (string, string) {
	log.Print("Generate Platform.sh structure...")

	pshDstAppPath := path.Join(rootPath, app, entity.PSH_PLATFORM)
	pshDstPath := path.Join(rootPath, PSH_CONFIG_PATH)
	return pshDstAppPath, pshDstPath
}

func TransfertConfigCustom(src string, dst string) {
	log.Println("Move custom config...")
	absProjectSourceConfig := filepath.Join(src, PSH_CONFIG_PATH)
	dirs, _ := ListDir(absProjectSourceConfig)
	for _, dir := range dirs {
		if dir != "local" {
			srcCustom := filepath.Join(absProjectSourceConfig, dir)
			dstCustom := filepath.Join(dst, dir)
			err := os.Mkdir(dstCustom, 0755)
			if err != nil {
				log.Println(err)
			}

			err = CopyDir(srcCustom, dstCustom)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

///// GENERIC SECTION /////

func (i *PathTmp) CleanUp() {
	log.Print("Clean up tempory workspace...")
	os.RemoveAll(i.Root)
}

func ResolvHomePath(path string) string {
	result := path

	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		result = filepath.Join(home, path[2:])
	}

	return result
}

func AbsPath(path string) string {
	absPath, _ := filepath.Abs(path)
	return absPath
}

func NormalizePath(path string) string {
	result := ResolvHomePath(path)
	result = AbsPath(result)

	return result
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func IsKeep() bool {
	return app.ArgsM.KeepData != ""
}

func CopyFile(src string, dst string) error {
	log.Printf("Copy file : %v to %v", src, dst)

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, errCopy := io.Copy(destination, source)

	err = os.Chmod(dst, sourceFileStat.Mode())
	if err != nil {
		return err
	}

	return errCopy
}

func ListDir(src string) ([]string, error) {
	var files []string
	f, err := os.Open(src)
	if err != nil {
		return files, err
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		if file.Mode().IsDir() {
			files = append(files, file.Name())
		}
	}
	return files, nil
}

func CopyDir(src string, dst string) error {
	var err error
	var fds []os.DirEntry
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if srcinfo.Mode().IsDir() {
		log.Printf("Copy Directory : %v to %v", src, dst)
		if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
			return err
		}

		if fds, err = os.ReadDir(src); err != nil {
			return err
		}
		for _, fd := range fds {
			srcfp := path.Join(src, fd.Name())
			dstfp := path.Join(dst, fd.Name())

			if fd.IsDir() {
				if err = CopyDir(srcfp, dstfp); err != nil {
					fmt.Println(err)
				}
			} else {
				if err = CopyFile(srcfp, dstfp); err != nil {
					fmt.Println(err)
				}
			}
		}
	} else {
		if err = CopyFile(src, dst); err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func GetFile(src string, dst string, ws string) {
	// Remote
	if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
		log.Printf("Get remote file from %v to %v ...", src, dst)

		response, err := http.Get(src)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer response.Body.Close()

		destination, err := os.Create(dst)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer destination.Close()

		_, err = io.Copy(destination, response.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// Local
	if strings.HasPrefix(src, "file://") {
		log.Printf("Get local file from %v to %v ...", src, dst)

		relPath := strings.Replace(src, "file://", "", 1)
		absPath := path.Join(ws, relPath)
		//absDest := path.Join(dst, relPath)
		err := CopyDir(absPath, dst)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func TimeStamp() string {
	ts := time.Now().UTC().Format(time.RFC3339)
	return strings.Replace(strings.Replace(ts, ":", "", -1), "-", "", -1)
}
