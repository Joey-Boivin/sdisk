package main

import (
	"fmt"
	"github.com/Joey-Boivin/sdisk-api/api/application"
	"github.com/Joey-Boivin/sdisk-api/api/repository"
	"log"
	"net/http"
	"os"

	"github.com/Joey-Boivin/sdisk-api/api/handlers"
	"gopkg.in/yaml.v3"
)

const (
	apiConfigPath = "./configs/api.yml"
)

type ApiConfig struct {
	Host string `yaml:"host"`
	Port uint   `yaml:"port"`
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
	userResource := handlers.NewUserHandler(registerService, fetchUserService)
	pingResource := handlers.NewPingHandler()

	router := http.NewServeMux()
	router.HandleFunc(handlers.PingEndpoint, pingResource.Get)
	router.HandleFunc(handlers.CreateUserEndpoint, userResource.Post)
	router.HandleFunc(handlers.GetUserEndpoint, userResource.Get)

	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", conf.Host, conf.Port), router); err != nil {
		log.Fatalf("Error trying to start server on %s:%d. %v", conf.Host, conf.Port, err)
	}
}
