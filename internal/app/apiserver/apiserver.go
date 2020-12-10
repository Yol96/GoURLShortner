package apiserver

import (
	"net/http"

	"github.com/Yol96/GoURLShortner/internal/app/store"
	"github.com/rs/cors"
)

// Start starts API server
func Start(config *Config) error {
	// Create a new store config
	db, err := store.NewConfig()
	if err != nil {
		return err
	}

	// Create a new storage
	store, err := store.NewStore(db)
	if err != nil {
		return err
	}

	defer store.Cli.Close()

	// Create a new configured server
	srv := newServer(store)
	srv.logger.Infof("Starting API server with next params: config:%+v db:%+v", config, db)
	handler := cors.Default().Handler(srv)
	
	return http.ListenAndServe(config.ServerPort, handler)
}
