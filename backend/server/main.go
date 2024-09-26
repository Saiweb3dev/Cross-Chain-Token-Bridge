package main

import (
	"flag"
	"log"
	"backend/server/mainserver"
	"backend/server/testserver"
)

func main() {
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
