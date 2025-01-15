package main

import "github.com/Joey-Boivin/sdisk/internal/infrastructure"

func main() {
	clientConfig := infrastructure.NewDefaultTCPClientConfig()
	client := infrastructure.NewTCPClient(clientConfig)
	client.Run()
}
