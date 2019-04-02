package main

import (
	"github.com/mimah/gomet/controller"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

func main() {

	rand.Seed(time.Now().UnixNano())

	err := os.MkdirAll("logs", 0700)
	if err != nil {
		fmt.Printf("Failed to create logs directory %s\n", err)
		return
	}

	err = os.MkdirAll("share", 0700)
	if err != nil {
		fmt.Printf("Failed to create share directory %s\n", err)
		return
	}

	logFile, _ := os.Create("logs/client.log")
	log.SetOutput(logFile)

	config, err := controller.LoadConfig()
	if err != nil {
		fmt.Printf("Invalid configuration file: %s\n", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)

	server := controller.NewServer(&wg, config)
	server.Start()

	cli := controller.NewCLI(server)
	go cli.Start()

	if config.Api.Enable {
		api := controller.NewApi(server)
		go api.Start()
	}

	log.Printf("Waiting for server to stop")
	wg.Wait()

	log.Printf("Server stopped")
}