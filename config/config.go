package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

var Conf config
var configName = "config.yml"
var configTemplateName = "config-example.yml"

type config struct {
	Port  string            `yaml:"port"`
	Auth  string            `yaml:"auth"`
	Files map[string]string `yaml:"files"`
}

func init() {
	if _, err := os.Stat(configName); os.IsNotExist(err) {

		workDir, _ := os.Getwd()
		GenerateFile(configTemplateName, filepath.Join(workDir), configName)
		fmt.Println("config.yml file created with default values.")
	}
	confFile, err := os.Open(configName)
	if err != nil {
		panic("Failed to open config.yml: " + err.Error())
	}
	defer confFile.Close()
	conf := config{}
	if err = yaml.NewDecoder(confFile).Decode(&conf); err != nil {
		panic(err)
	}
	Conf = conf
}
