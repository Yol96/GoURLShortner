package store

import (
	"github.com/go-redis/redis"
)

// Store contains storage struct with pointers to Redis Client, User Repository
type Store struct {
	Cli            *redis.Client
	userRepository *UserRepository
}

// User creates a new userrepository
func (cli *Store) User() *UserRepository {
	// TODO: add test store for testing
	if cli.userRepository != nil {
		return cli.userRepository
	}

	cli.userRepository = &UserRepository{
		store: cli,
	}

	return cli.userRepository
}

// NewStore creates a new configured store from config
func NewStore(config *Config) (*Store, error) {
	config, err := NewConfig()
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     config.address,
		Password: config.password,
		DB:       config.db,
	})

	// Checking connection to db
	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	store := &Store{
		Cli: client,
	}

	return store, nil
}
