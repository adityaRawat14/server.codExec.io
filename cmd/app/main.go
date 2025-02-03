package main

import (
	"log"
	"os"
	"os/signal"
	"server/internal/api/http/routes"
	db "server/internal/database"
	"server/services/util/executor"
	"server/services/util/k8s"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	err := executor.GetNewExecutorClient()
	if err != nil {
		log.Println("Failed to initiate docker client!!")
		os.Exit(1)
	}
	log.Println("Docker client initialized!")
	_,err= k8s.NewClient()
	if err != nil {
		log.Println("Failed to initiate kubernetes client!!")
		os.Exit(1)
	}
	router := gin.New()

	// Register all the routes
	router.Use(cors.Default())
	routes.UserRoutes(router)
	routes.CodeRoutes(router)
	// setting up cors
	// Open database pool
	err = db.OpenDbPool()
	if err != nil {
		log.Println("Failed to create database pool !!")
		os.Exit(1)
	}

	// Channel to listen for OS events
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := router.Run(":8000"); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()

	sig := <-signalChan
	log.Printf("Received signal: %s. Shutting down gracefully...", sig)

	// Shutting down all the listeners
	log.Println("Closing the DB connection pool !!")
	db.ShutDownDbPool()
	// log.Println("Closing the Docker client")
	// executor.CloseClient()

	log.Println("Server shut down gracefully")
}
