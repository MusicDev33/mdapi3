package main

import (
	"MusicDev33/mdapi3/internal/server"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Port int `yaml:"port"`
}

func main() {
	f, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	var cfg Config
	if err := yaml.Unmarshal(f, &cfg); err != nil {
		panic(err)
	}

	srv := server.NewServer()

	fmt.Printf("Starting server on port %d...\n", cfg.Port)
	if err := srv.Run(cfg.Port); err != nil {
		panic(err)
	}
}
