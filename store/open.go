package store

import (
	"database/sql"
	"github.com/bmizerany/pq"
	"os"
)

func Open(url string) (*sql.DB, error) {
	if url == "" {
		url = os.Getenv("DATABASE_URL")
	}

	if url == "" {
		url = os.Getenv("BOXEN_POSTGRESQL_URL")
		if url != "" {
			url += "simplex?sslmode=disable"
		}
	}

	src, err := pq.ParseURL(url)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", src)
	if err != nil {
		return nil, err
	}

	return db, nil
}
