package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	dbHost     = "localhost"
	dbPort     = "5432"
	dbUser     = "postgres"
	dbPassword = "Manish@2001"
	dbName     = "Catalog"
)

func NewDBConnection() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	err = createTables(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS restaurants (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			location VARCHAR(255) NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS menu_items (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			price DECIMAL(10,2) NOT NULL,
			restaurant_id INTEGER REFERENCES restaurants(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return err
	}

	return nil
}
