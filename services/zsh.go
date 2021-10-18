package services

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"text/template"
)

var SystemOptions = []string{"arch", "osx", "ubuntu"}

type ParseOptions struct {
	IsOsx       bool
	IsUbuntu    bool
	IsArch      bool
	SupportSnap bool
	SupportJava bool
}

func CreateZshFiles(dest string, system string, supportJava bool) error {
	templateParser := template.New("parser")
	parseOptions := makeParseOptions(system, supportJava)
	tmplDir := path.Join(rootPath(), "templates/")
	templateFiles, err := os.ReadDir(tmplDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range templateFiles {
		tmplFile := path.Join(tmplDir, file.Name())
		tmpl, err := templateParser.ParseFiles(tmplFile)
		if err != nil {
			return err
		}

		destFilename := fmt.Sprintf("%s/.%s", dest, file.Name())
		destFileWriter, err := os.Create(destFilename)
		if err != nil {
			log.Fatal(err)
		}

		defer func() {
			if err := destFileWriter.Close(); err != nil {
				panic(err)
			}
		}()

		if err := tmpl.Execute(destFileWriter, parseOptions); err != nil {
			return err
		}
	}

	return nil
}

func IsSystemOption(value string) bool {
	for _, op := range SystemOptions {
		if op == value {
			return true
		}
	}
	return false
}

func makeParseOptions(system string, supportJava bool) ParseOptions {
	return ParseOptions{
		IsArch:      system == "arch",
		IsOsx:       system == "osx",
		IsUbuntu:    system == "ubuntu",
		SupportSnap: system == "ubuntu" || system == "arch",
		SupportJava: supportJava,
	}
}

func rootPath() string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Join(filepath.Dir(b), "..")
	return basepath
}
