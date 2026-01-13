package main

import (
	"exc8/client"
	"exc8/server"
	"log"
	"time"
)

func main() {
	// Start gRPC server in a goroutine
	go func() {
		if err := server.StartGrpcServer(); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()
	// Give him a second to start
	time.Sleep(1 * time.Second)

	// Create gRPC client and run
	c, err := client.NewGrpcClient()
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	if err := c.Run(); err != nil {
		log.Fatalf("failed to run client: %v", err)
	}

	println("Orders complete!")
}
