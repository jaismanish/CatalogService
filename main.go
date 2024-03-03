package main

import (
	"database/sql"
	"fmt"
	"github.com/jaismanish15/CatalogService/server"
	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() error {
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

func main() {
	err := initDB()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}

	defer db.Close()

	// Create an instance of the catalog service
	catalogService := catalog.NewCatalogService(db)

	// Example usage
	restaurant, err := catalogService.AddRestaurant("New Restaurant", "New Location")
	if err != nil {
		fmt.Println("Error adding restaurant:", err)
		return
	}
	fmt.Println("Added restaurant:", restaurant)

	menuItem, err := catalogService.AddMenuItem("New Item", 12.99, restaurant.Id)
	if err != nil {
		fmt.Println("Error adding menu item:", err)
		return
	}
	fmt.Println("Added menu item:", menuItem)

	restaurants, err := catalogService.GetRestaurants()
	if err != nil {
		fmt.Println("Error getting restaurants:", err)
		return
	}
	fmt.Println("Restaurants:", restaurants)

	menuItems, err := catalogService.GetMenuItems()
	if err != nil {
		fmt.Println("Error getting menu items:", err)
		return
	}
	fmt.Println("Menu Items:", menuItems)
}
