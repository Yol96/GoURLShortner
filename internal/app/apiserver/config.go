package apiserver

import (
	"os"

	"github.com/joho/godotenv"
)

// Config contains server config parameters(port, logrus log level)
type Config struct {
	ServerPort  string
	LogrusLevel string
}

// NewConfig creates a new application config
func NewConfig() (*Config, error) {
	// Getting config fields from .env
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	address := os.Getenv("APP_SERVER_PORT")
	if address == "" {
		address = ":8080"
	}

	logLevel := os.Getenv("APP_LOGRUS_LEVEL")
	if logLevel == "" {
		logLevel = "debug"
	}

	return &Config{
		ServerPort:  address,
		LogrusLevel: logLevel,
	}, nil
}
