package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const (
	createUserTable = `
	CREATE TABLE IF NOT EXISTS user (
		id VARCHAR(255) NOT NULL,
		username VARCHAR(255) NOT NULL,
	    email VARCHAR(255) NOT NULL,
		password VARCHAR(256) NOT NULL,
		PRIMARY KEY (id)
	);`
)

func NewDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./instance/mlanal.db")

	if err != nil {
		return nil, err
	}
	
	if _, err := db.Exec(createUserTable); err != nil {
		return nil, err
	}

	return db, nil
}
