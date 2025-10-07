package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DbConfig struct {
	Host          string
	User          string
	Port          int
	Name          string
	Password      string
	EnableSSLMode bool
}

type Config struct {
	Version     string
	HttpPort    string
	ServiceName string
	Secret      string
	Db          DbConfig
}

var configurations *Config

func loadConfig() {

	err := godotenv.Overload()
	if err != nil {
		fmt.Println("Failed to load .env file")
		os.Exit(1)
	}

	version := os.Getenv("VERSION")
	httpPort := os.Getenv("HTTP_PORT")
	serviceName := os.Getenv("SERVICE_NAME")
	secret := os.Getenv("SECRET")

	if version == "" || httpPort == "" || serviceName == "" || secret == "" {
		fmt.Println("Missing required app env variables")
		os.Exit(1)
	}

	host := os.Getenv("HOST")
	user := os.Getenv("USER")
	portStr := os.Getenv("PORT")
	dbName := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")
	enableSSLStr := os.Getenv("ENABLE_SSL_MODE")

	if host == "" || user == "" || portStr == "" || dbName == "" || password == "" || enableSSLStr == "" {
		fmt.Println("Missing required DB env variables")
		os.Exit(1)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		fmt.Println("Invalid PORT value in .env")
		os.Exit(1)
	}

	enableSSL, err := strconv.ParseBool(enableSSLStr)
	if err != nil {
		fmt.Println("Invalid ENABLE_SSL_MODE value in .env")
		os.Exit(1)
	}

	configurations = &Config{
		Version:     version,
		HttpPort:    httpPort,
		ServiceName: serviceName,
		Secret:      secret,
		Db: DbConfig{
			Host:          host,
			User:          user,
			Port:          port,
			Name:          dbName,
			Password:      password,
			EnableSSLMode: enableSSL,
		},
	}
}

func Load() *Config {
	if configurations == nil {
		loadConfig()
	}
	return configurations
}
