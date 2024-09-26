package main

import (
	"flag"
	"log"
	"backend/config"
	"backend/server/mainserver"
	"backend/server/testserver"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
}
	serverType := flag.String("server", "main", "Specify which server to start (main or test)")
	flag.Parse()

	switch *serverType {
	case "main":
		mainserver.RunMainServer()
	case "test":
		testserver.RunTestServer()
	default:
		log.Fatal("Invalid server type specified")
	}
}
