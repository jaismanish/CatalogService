package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() error {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "postgres", "Manish@2001", "Catalog")

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	err = createTables()
	if err != nil {
		return err
	}

	return nil
}

func createTables() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS restaurants (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255),
			location VARCHAR(255)
		);
		CREATE TABLE IF NOT EXISTS menu_items (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255),
			price DOUBLE PRECISION,
			restaurant_id INT
		);
	`)

	return err
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}

func GetDB() *sql.DB {
	return db
}
