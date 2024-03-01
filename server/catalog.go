package main

import (
	"context"
	"sync"
)

type CatalogServiceServer struct {
	mu           sync.Mutex
	restaurants  map[int32]*Restaurant
	menuItems    map[int32][]*MenuItem
	restaurantID int32
	menuItemID   int32
}

type Restaurant struct {
	Id       int32
	Name     string
	Location string
}

type MenuItem struct {
	Id           int32
	Name         string
	Price        float64
	RestaurantId int32
}

func (s *CatalogServiceServer) AddRestaurant(ctx context.Context, req *Restaurant) (*Restaurant, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.restaurantID++
	restaurant := &Restaurant{
		Id:       s.restaurantID,
		Name:     req.Name,
		Location: req.Location,
	}

	s.restaurants[s.restaurantID] = restaurant

	return restaurant, nil
}

func (s *CatalogServiceServer) AddMenuItem(ctx context.Context, req *MenuItem) (*MenuItem, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.menuItemID++
	menuItem := &MenuItem{
		Id:           s.menuItemID,
		Name:         req.Name,
		Price:        req.Price,
		RestaurantId: req.RestaurantId,
	}

	s.menuItems[req.RestaurantId] = append(s.menuItems[req.RestaurantId], menuItem)

	return menuItem, nil
}
