package main

import (
	"fmt"
	"github.com/jaismanish15/CatalogService/db"
	"github.com/jaismanish15/CatalogService/server"
	_ "github.com/lib/pq"
)

func main() {
	err := db.InitDB()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}

	defer db.CloseDB()

	catalogService := catalog.NewCatalogService(db.GetDB())

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
