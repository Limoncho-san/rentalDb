package rental

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// InitDB initializes the database connection

func InitDB(connectionString string) error {
	var err error
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	log.Println("Database connection successful")
	return nil
}

// CloseDB closes the database connection

func CloseDB() {
	db.Close()
}
