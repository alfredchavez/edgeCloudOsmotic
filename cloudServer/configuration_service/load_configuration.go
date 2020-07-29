package configuration_service

import (
	"github.com/goccy/go-yaml"
	"io/ioutil"
	"log"
	"sync"
)

var loadOnceConfig sync.Once

var mainServerAddr string
var externalServerName string
var internalPort string
var runtime string

type Configuration struct {
	Runtime string `yaml:"runtime"`
	ServerConfiguration struct{
		Port string `yaml:"port"`
		Name string `yaml:"name"`
	} `yaml:"server_configuration"`
	MainServerConfiguration struct {
		Address string `yaml:"addr"`
	} `yaml:"main_server_configuration"`
}

func LoadConfiguration(){
	loadOnceConfig.Do(func() {
		mainServerAddr = "http://localhost:9000"
		externalServerName = "edge_1"
		internalPort = "8000"
		runtime = "wasmer"

		yamlFile, err := ioutil.ReadFile("config.yml")
		if err != nil {
			log.Printf("Could not read the file(ReadFile) %v", err)
			return
		}
		config := Configuration{}
		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			log.Printf("Could not assign variables to struct(Unmarshal) %v", err)
			return
		}
		mainServerAddr = config.MainServerConfiguration.Address
		externalServerName = config.ServerConfiguration.Name
		internalPort = config.ServerConfiguration.Port
		runtime = config.Runtime
	})
}

func GetMainServerAddr() string {
	return mainServerAddr
}

func GetMyServerName() string {
	return externalServerName
}

func GetMyServerPort() string {
	return internalPort
}

func GetRuntime() string {
	return runtime
}

