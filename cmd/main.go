package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Joey-Boivin/sdisk/internal/application"
	"github.com/Joey-Boivin/sdisk/internal/handlers"
	"github.com/Joey-Boivin/sdisk/internal/repository"
	"gopkg.in/yaml.v3"
)

const (
	apiConfigPath = "./configs/api.yml"
)

type ApiConfig struct {
	Host     string `yaml:"host"`
	Port     uint   `yaml:"port"`
	DiskSize uint   `yaml:"diskSize"`
}

func main() {
	file, err := os.Open(apiConfigPath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	defer file.Close()

	var conf ApiConfig
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&conf); err != nil {
		log.Fatalf("Error decoding yaml file: %v", err)
	}

	log.Printf("Server starting on %s:%d", conf.Host, conf.Port)

	userRepository := repository.NewRamRepository()
	registerService := application.NewRegisterService(userRepository)
	fetchUserService := application.NewFetchUserService(userRepository)
	createDiskService := application.NewCreateDiskService(userRepository, uint64(conf.DiskSize))
	userResource := handlers.NewUserHandler(registerService, fetchUserService, createDiskService)
	pingResource := handlers.NewPingHandler()

	router := http.NewServeMux()
	router.HandleFunc(handlers.PingEndpoint, pingResource.Ping)
	router.HandleFunc(handlers.CreateUserEndpoint, userResource.CreateUserResource)
	router.HandleFunc(handlers.GetUserEndpoint, userResource.GetUserResource)
	router.HandleFunc(handlers.CreateDiskEndpoint, userResource.CreateDiskResource)

	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", conf.Host, conf.Port), router); err != nil {
		log.Fatalf("Error trying to start server on %s:%d. %v", conf.Host, conf.Port, err)
	}
}
