package main

import (
	"fmt"
	"os"

	"github.com/Joey-Boivin/sdisk/internal/infrastructure"
	"github.com/Joey-Boivin/sdisk/internal/models"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a token.")
		return
	}

	id := os.Args[1]
	userID, err := models.FromString(id)
	if err != nil {
		panic(err)
	}

	clientConfig := infrastructure.NewDefaultTCPClientConfig(userID)
	client := infrastructure.NewTCPClient(clientConfig)
	client.Run()
}
