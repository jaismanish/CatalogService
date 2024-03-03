package catalog

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testDB *sql.DB

func setupTestDB() error {
	connStr := "host=localhost port=5432 user=postgres password=Manish@2001 dbname=Catalog sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
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
	if err != nil {
		return err
	}

	testDB = db
	return nil
}

func teardownTestDB() {
	testDB.Exec("DROP TABLE IF EXISTS menu_items CASCADE")
	testDB.Exec("DROP TABLE IF EXISTS restaurants CASCADE")
	testDB.Close()
}

func TestMain(m *testing.M) {
	err := setupTestDB()
	if err != nil {
		panic(err)
	}

	code := m.Run()

	teardownTestDB()

	os.Exit(code)
}

func TestCatalogService_AddRestaurant(t *testing.T) {
	catalogService := NewCatalogService(testDB)

	restaurant, err := catalogService.AddRestaurant("Test Restaurant", "Test Location")
	if err != nil {
		t.Fatalf("Error adding restaurant: %v", err)
	}

	if restaurant.Id == 0 {
		t.Error("Expected non-zero ID for the added restaurant")
	}

	if restaurant.Name != "Test Restaurant" {
		t.Errorf("Expected restaurant name to be 'Test Restaurant', got '%s'", restaurant.Name)
	}

	if restaurant.Location != "Test Location" {
		t.Errorf("Expected restaurant location to be 'Test Location', got '%s'", restaurant.Location)
	}
}

func TestCatalogService_AddMenuItem(t *testing.T) {
	catalogService := NewCatalogService(testDB)

	restaurant, err := catalogService.AddRestaurant("Test Restaurant", "Test Location")
	if err != nil {
		t.Fatalf("Error adding restaurant: %v", err)
	}

	menuItem, err := catalogService.AddMenuItem("Test Item", 9.99, restaurant.Id)
	if err != nil {
		t.Fatalf("Error adding menu item: %v", err)
	}

	if menuItem.Id == 0 {
		t.Error("Expected non-zero ID for the added menu item")
	}

	if menuItem.Name != "Test Item" {
		t.Errorf("Expected menu item name to be 'Test Item', got '%s'", menuItem.Name)
	}

	if menuItem.Price != 9.99 {
		t.Errorf("Expected menu item price to be 9.99, got %.2f", menuItem.Price)
	}

	if menuItem.RestaurantId != restaurant.Id {
		t.Errorf("Expected menu item restaurant ID to be %d, got %d", restaurant.Id, menuItem.RestaurantId)
	}
}

func TestCatalogService_GetRestaurants(t *testing.T) {
	catalogService := NewCatalogService(testDB)

	_, err := catalogService.AddRestaurant("Restaurant 1", "Location 1")
	if err != nil {
		t.Fatalf("Error adding restaurant: %v", err)
	}

	_, err = catalogService.AddRestaurant("Restaurant 2", "Location 2")
	if err != nil {
		t.Fatalf("Error adding restaurant: %v", err)
	}

	restaurants, err := catalogService.GetRestaurants()
	if err != nil {
		t.Fatalf("Error getting restaurants: %v", err)
	}

	if len(restaurants) != 2 {
		t.Errorf("Expected 2 restaurants, got %d", len(restaurants))
	}
}

func TestCatalogService_GetMenuItems(t *testing.T) {
	catalogService := NewCatalogService(testDB)

	restaurant, err := catalogService.AddRestaurant("Test Restaurant", "Test Location")
	if err != nil {
		t.Fatalf("Error adding restaurant: %v", err)
	}

	_, err = catalogService.AddMenuItem("Item 1", 10.99, restaurant.Id)
	if err != nil {
		t.Fatalf("Error adding menu item: %v", err)
	}

	_, err = catalogService.AddMenuItem("Item 2", 15.99, restaurant.Id)
	if err != nil {
		t.Fatalf("Error adding menu item: %v", err)
	}

	menuItems, err := catalogService.GetMenuItems()
	if err != nil {
		t.Fatalf("Error getting menu items: %v", err)
	}

	if len(menuItems) != 2 {
		t.Errorf("Expected 2 menu items, got %d", len(menuItems))
	}
}
