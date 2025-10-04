package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Version     string
	HttpPort    string
	ServiceName string
	Secret      string
}

var configurations *Config

func loadConfig() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load env")
		os.Exit(1)
	}

	version := os.Getenv("VERSION")
	httpPort := os.Getenv("HTTP_PORT")
	serviceName := os.Getenv("SERVICE_NAME")
	secret := os.Getenv("SECRET")

	if version == "" || httpPort == "" || serviceName == "" || secret == "" {
		fmt.Println("Missing required env variables")
		os.Exit(1)
	}

	// ✅ Initialize struct before assigning
	configurations = &Config{
		Version:     version,
		HttpPort:    httpPort,
		ServiceName: serviceName,
		Secret:      secret,
	}
}

func Load() *Config {
	if configurations == nil {
		loadConfig()
		return configurations
	} else {
		return configurations
	}
}
