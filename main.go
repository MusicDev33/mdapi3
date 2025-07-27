package main

import (
	"MusicDev33/mdapi3/internal/config"
	"MusicDev33/mdapi3/internal/server"
	"context"
	"fmt"
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

	uri := fmt.Sprintf("mongodb://localhost:%d", cfg.MongoPort)

	_, err := InitMongo(uri, cfg.DBName)
	if err != nil {
		panic(err)
	}

	srv := server.NewServer()

	fmt.Printf("Starting server on port %d...\n", cfg.Port)
	if err := srv.Run(cfg.Port); err != nil {
		panic(err)
	}
}
