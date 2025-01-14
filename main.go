package main

import (
	"app/config"
	"app/server"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	log.Println("Starting Runners App")
	log.Println("Initializing configuration")
	config := config.InitConfig("runners")
	log.Println("Initializing database")
	dbHandler := server.InitDatabase(config)
	log.Println("Initializing HTTP server")
	httpServer := server.InitHttpServer(config, dbHandler)
	httpServer.Start()
}
