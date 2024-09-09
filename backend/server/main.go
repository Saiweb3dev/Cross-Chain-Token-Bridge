package main

import (
	"backend/routes"
	"backend/config"
)

func main() {
	r := routes.SetupRouter()
	r.Run(config.ServerAddress())
}



