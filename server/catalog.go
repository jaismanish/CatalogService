package catalog

import (
	"database/sql"
)

type Restaurant struct {
	Id       int32  `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

type MenuItem struct {
	Id           int32   `json:"id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	RestaurantId int32   `json:"restaurantId"`
}

type CatalogService interface {
	AddRestaurant(name, location string) (*Restaurant, error)
	AddMenuItem(name string, price float64, restaurantID int32) (*MenuItem, error)
	GetRestaurants() ([]*Restaurant, error)
	GetMenuItems() ([]*MenuItem, error)
}

func NewCatalogService(db *sql.DB) CatalogService {
	return &catalogService{db: db}
}

type catalogService struct {
	db *sql.DB
}

func (c *catalogService) AddRestaurant(name, location string) (*Restaurant, error) {
	query := "INSERT INTO restaurants (name, location) VALUES ($1, $2) RETURNING id"
	var id int32
	err := c.db.QueryRow(query, name, location).Scan(&id)
	if err != nil {
		return nil, err
	}

	restaurant := &Restaurant{
		Id:       id,
		Name:     name,
		Location: location,
	}

	return restaurant, nil
}

func (c *catalogService) AddMenuItem(name string, price float64, restaurantID int32) (*MenuItem, error) {
	query := "INSERT INTO menu_items (name, price, restaurant_id) VALUES ($1, $2, $3) RETURNING id"
	var id int32
	err := c.db.QueryRow(query, name, price, restaurantID).Scan(&id)
	if err != nil {
		return nil, err
	}

	menuItem := &MenuItem{
		Id:           id,
		Name:         name,
		Price:        price,
		RestaurantId: restaurantID,
	}

	return menuItem, nil
}

func (c *catalogService) GetRestaurants() ([]*Restaurant, error) {
	rows, err := c.db.Query("SELECT id, name, location FROM restaurants")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var restaurants []*Restaurant
	for rows.Next() {
		restaurant := &Restaurant{}
		err := rows.Scan(&restaurant.Id, &restaurant.Name, &restaurant.Location)
		if err != nil {
			return nil, err
		}
		restaurants = append(restaurants, restaurant)
	}

	return restaurants, nil
}

func (c *catalogService) GetMenuItems() ([]*MenuItem, error) {
	rows, err := c.db.Query("SELECT id, name, price, restaurant_id FROM menu_items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var menuItems []*MenuItem
	for rows.Next() {
		menuItem := &MenuItem{}
		err := rows.Scan(&menuItem.Id, &menuItem.Name, &menuItem.Price, &menuItem.RestaurantId)
		if err != nil {
			return nil, err
		}
		menuItems = append(menuItems, menuItem)
	}

	return menuItems, nil
}
