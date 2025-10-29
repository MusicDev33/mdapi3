package main

import (
	"MusicDev33/mdapi3/internal/config"
	"MusicDev33/mdapi3/internal/database"
	"MusicDev33/mdapi3/internal/server"
	"fmt"
	"log"
)

func main() {
	cfg := config.Get()

	// Initialize MongoDB connection
	uri := fmt.Sprintf("mongodb://localhost:%d", cfg.MongoPort)
	db, err := database.InitMongo(uri, cfg.DBName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.Close()

	// Set global database instance
	database.DB = db

	fmt.Printf("Connected to MongoDB at %s\n", uri)

	// Create and start server
	srv := server.NewServer()

	fmt.Printf("Starting server on port %d...\n", cfg.Port)
	if err := srv.Run(cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
