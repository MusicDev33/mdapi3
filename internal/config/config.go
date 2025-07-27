package config

import (
	"os"
	"sync"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Port          int    `yaml:"port"`
	MongoPort     int    `yaml:"mongoPort"`
	DBName        string `yaml:"dbName"`
	KeyAnthropic  string `yaml:"akAnthropic"`
	KeyDeepSeek   string `yaml:"akDeepSeek"`
	KeyOpenAI     string `yaml:"akOpenAI"`
	WhitelistCORS string `yaml:"whitelistCors"`
}

var (
	instance *Config
	once     sync.Once
)

func Get() *Config {
	once.Do(func() {
		f, err := os.ReadFile("config.yaml")
		if err != nil {
			panic(err)
		}

		var cfg Config
		if err := yaml.Unmarshal(f, &cfg); err != nil {
			panic(err)
		}

		instance = &cfg
	})

	return instance
}
