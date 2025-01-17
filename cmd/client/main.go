package main

import (
	"log"
	"os"

	"github.com/Joey-Boivin/sdisk/internal/infrastructure"
	"github.com/Joey-Boivin/sdisk/internal/models"
	"gopkg.in/yaml.v3"
)

type ClientConfig struct {
	Host       string `yaml:"host"`
	Port       uint   `yaml:"port"`
	FolderName string `yaml:"folderName"`
	Token      string `yaml:"token"`
}

func main() {
	path := os.Getenv("SDISK_HOME")
	path += "/configs/client.yml"

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	defer file.Close()

	var conf ClientConfig
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&conf); err != nil {
		log.Fatalf("Error decoding yaml file: %v", err)
	}

	userID, err := models.FromString(conf.Token)

	if err != nil {
		panic(err)
	}

	clientConfig := infrastructure.NewDefaultTCPClientConfig(userID, conf.Host, conf.Port, conf.FolderName)
	client := infrastructure.NewTCPClient(clientConfig)
	client.Run()
}
