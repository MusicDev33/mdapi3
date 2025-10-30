package main

import (
	"MusicDev33/mdapi3/internal/config"
	"MusicDev33/mdapi3/internal/database"
	"MusicDev33/mdapi3/internal/server"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func InitMongo(uri string, dbName string) (*Mongo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		panic(err)
	}

	db := client.Database(dbName)
	return &Mongo{Client: client, DB: db}, nil
}

func (m *Mongo) Close() error {
	return m.Client.Disconnect(context.Background())
}

func main() {
	cfg := config.Get()

	// Initialize MongoDB connection
	uri := fmt.Sprintf("mongodb://%s:%d", cfg.MongoURI, cfg.MongoPort)
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
