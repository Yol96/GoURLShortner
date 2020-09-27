package store

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config contains storage config parameters(address, password, db)
type Config struct {
	address  string
	password string
	db       int
}

// NewConfig creates a new storage config
func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	address := os.Getenv("APP_REDIS_ADDRESS")
	if address == "" {
		address = "localhost:6379"
	}

	password := os.Getenv("APP_REDIS_PASSWORD")
	if password == "" {
		password = ""
	}

	db := os.Getenv("APP_REDIS_DB")
	if db == "" {
		db = "0"
	}

	index, err := strconv.Atoi(db)
	if err != nil {
		return nil, err
	}

	return &Config{
		address:  address,
		password: password,
		db:       index,
	}, nil
}
