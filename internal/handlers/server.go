package handlers

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/theo-mazars/postrun/internal/config"
)

type Server struct {
	cfg  *config.Config
	pool *sql.DB
}

func Load() Server {
	cfg := config.Load("")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Pass,
		cfg.Database.Name,
	)
	pool, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("unable to use data source name", err)
	}

	err = pool.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return Server{cfg, pool}
}
