package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

type Version struct {
	Latest   VersionInfo   `json:"latest"`
	Versions []VersionInfo `json:"versions"`
}

type VersionInfo struct {
	Date string `json:"date"`
	Sha  string `json:"sha"`
	Name string `json:"name"`
}

type Definition struct {
	Names    []string  `json:"names"`
	Commands []command `json:"commands"`
}

type command struct {
	Name    string      `json:"name"`
	Usage   string      `json:"usage"`
	Type    commandType `json:"type"`
	Command string      `json:"command"`
}

type commandType string

const (
	shell  commandType = "shell"
	docker             = "docker"
	url                = "url"
)

func main() {
	app := cli.NewApp()
	app.Name = "dockleaf"
	app.Usage = "Ever changing dev/ops functions, in a consistent way..."
	app.Action = func(c *cli.Context) error {

		definition, version := getInputs(c.Args())

		fmt.Println(definition)
		fmt.Println(version)
		cmd := exec.Command("go", "build", "-o", definition.Names[0])
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
			return nil
		}
		fmt.Println("Result: " + out.String())

		createOtherNames(definition.Names)

		return nil
	}

	app.Run(os.Args)
}

func createOtherNames(names []string) {
	path := os.Getenv("PWD")
	target := names[0]
	for i, name := range names {
		if i != 0 {
			symlink := filepath.Join(path, name)
			os.Symlink(target, symlink)
		}
	}
}

func getInputs(args cli.Args) (Definition, Version) {

	var definitionFile, versionFile string

	if args.Present() {
		definitionFile = args.Get(0)
		versionFile = args.Get(1)
	} else {
		definitionFile = os.Getenv("DOCKLEAF_DEFINITION")
		versionFile = os.Getenv("DOCKLEAF_VERSION")
	}

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		fmt.Println(pair[0], "=", pair[1])
	}

	fmt.Println("Definition:" + definitionFile)
	fmt.Println("Version:" + versionFile)

	definition := toDefinition(definitionFile)
	version := toVersion(versionFile)
	return definition, version
}

func readFile(filename string) []byte {
	filecontents, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Couldn't find the file: " + filename)
		} else {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
	return filecontents
}

func toDefinition(filename string) Definition {
	var definition Definition
	json.Unmarshal(readFile(filename), &definition)
	return definition
}

func toVersion(filename string) Version {
	var version Version
	if len(filename) > 0 {
		json.Unmarshal(readFile(filename), &version)
	} else {
		version = Version{}
	}

	return version
}
