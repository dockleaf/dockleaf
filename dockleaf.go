package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

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

		definitionFile := c.Args().Get(0)
		definition := toDefinition(definitionFile)
		var version Version
		if len(c.Args()) > 1 {
			versionFile := c.Args().Get(1)
			version = toVersion(versionFile)
		} else {
			version = Version{}
		}

		fmt.Println(definition)
		fmt.Println(version)
		cmd := exec.Command("go", "build", "-o", definition.Names[0])
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		return nil
	}

	app.Run(os.Args)
}

func readFile(filename string) []byte {
	filecontents, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("couldn't find the file")
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
	json.Unmarshal(readFile(filename), &version)
	return version
}