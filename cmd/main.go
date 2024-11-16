package main

import (
	"fmt"
	"github.com/Joey-Boivin/sdisk-api/api/application"
	"github.com/Joey-Boivin/sdisk-api/api/repository"
	"log"
	"net/http"
	"os"

	"github.com/Joey-Boivin/sdisk-api/api/handlers"
	"github.com/Joey-Boivin/sdisk-api/api/router"
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
	userResource := handlers.NewUserHandler(registerService)
	pingResource := handlers.NewPingHandler()

	router := router.NewRouter()
	router.AddRoute(pingResource.Get, http.MethodGet, handlers.PingEndpoint)
	router.AddRoute(userResource.Post, http.MethodPost, handlers.UsersEndpoint)

	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", conf.Host, conf.Port), router); err != nil {
		log.Fatalf("Error trying to start server on %s:%d. %v", conf.Host, conf.Port, err)
	}
}
