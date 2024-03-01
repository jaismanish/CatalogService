package main

import (
	"context"
	"testing"
)

func TestCatalogServiceServer_AddRestaurant(t *testing.T) {
	s := &CatalogServiceServer{
		restaurants: make(map[int32]*Restaurant),
		menuItems:   make(map[int32][]*MenuItem),
	}

	req := &Restaurant{
		Name:     "Restaurant1",
		Location: "Location1",
	}

	res, err := s.AddRestaurant(context.Background(), req)
	if err != nil {
		t.Fatalf("AddRestaurant failed: %v", err)
	}

	if res.Id != 1 {
		t.Errorf("Expected Restaurant ID to be 1, got %d", res.Id)
	}
}

func TestCatalogServiceServer_AddMenuItem(t *testing.T) {
	s := &CatalogServiceServer{
		restaurants:  make(map[int32]*Restaurant),
		menuItems:    make(map[int32][]*MenuItem),
		restaurantID: 1,
	}

	req := &MenuItem{
		Name:         "Test Item",
		Price:        9.99,
		RestaurantId: 1,
	}

	res, err := s.AddMenuItem(context.Background(), req)
	if err != nil {
		t.Fatalf("AddMenuItem failed: %v", err)
	}

	if res.Id != 1 {
		t.Errorf("Expected MenuItem ID to be 1, got %d", res.Id)
	}
}
