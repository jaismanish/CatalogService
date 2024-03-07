package main

import (
	"CatalogService/proto"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type catalogService struct {
	db *sql.DB
	proto.UnimplementedCatalogServiceServer
}

func NewCatalogService() (*catalogService, error) {
	db, err := NewDBConnection()
	if err != nil {
		return nil, err
	}

	return &catalogService{
		db: db,
	}, nil
}

func (s *catalogService) AddRestaurant(ctx context.Context, req *proto.AddRestaurantRequest) (*proto.AddRestaurantResponse, error) {
	if s.db == nil {
		return nil, errors.New("database connection is nil")
	}
	_, err := s.db.Exec("INSERT INTO restaurants (name, location) VALUES ($1, $2)", req.Name, req.Location)
	if err != nil {
		return nil, fmt.Errorf("failed to add restaurant: %v", err)
	}

	return &proto.AddRestaurantResponse{Success: true}, nil
}

func (s *catalogService) AddMenuItem(ctx context.Context, req *proto.AddMenuItemRequest) (*proto.AddMenuItemResponse, error) {
	_, err := s.db.Exec("INSERT INTO menu_items (name, price, restaurant_id) VALUES ($1, $2, $3)", req.Name, req.Price, req.RestaurantId)
	if err != nil {
		return nil, fmt.Errorf("failed to add menu item: %v", err)
	}

	return &proto.AddMenuItemResponse{Success: true}, nil
}

func (s *catalogService) GetRestaurants(ctx context.Context, req *proto.GetRestaurantsRequest) (*proto.GetRestaurantsResponse, error) {
	rows, err := s.db.Query("SELECT id, name, location FROM restaurants")
	if err != nil {
		return nil, fmt.Errorf("failed to get restaurants: %v", err)
	}
	defer rows.Close()

	var restaurants []*proto.Restaurant
	for rows.Next() {
		var r proto.Restaurant
		if err := rows.Scan(&r.Id, &r.Name, &r.Location); err != nil {
			return nil, fmt.Errorf("failed to scan restaurant: %v", err)
		}
		restaurants = append(restaurants, &r)
	}

	return &proto.GetRestaurantsResponse{Restaurants: restaurants}, nil
}

func (s *catalogService) GetMenuItems(ctx context.Context, req *proto.GetMenuItemsRequest) (*proto.GetMenuItemsResponse, error) {
	rows, err := s.db.Query("SELECT id, name, price, restaurant_id FROM menu_items WHERE restaurant_id = $1", req.RestaurantId)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu items: %v", err)
	}
	defer rows.Close()

	var menuItems []*proto.MenuItem
	for rows.Next() {
		var mi proto.MenuItem
		if err := rows.Scan(&mi.Id, &mi.Name, &mi.Price, &mi.RestaurantId); err != nil {
			return nil, fmt.Errorf("failed to scan menu item: %v", err)
		}
		menuItems = append(menuItems, &mi)
	}
	return &proto.GetMenuItemsResponse{MenuItems: menuItems}, nil
}

func main() {
	dbService, err := NewCatalogService()
	if err != nil {
		log.Fatalf("Failed to create catalog service: %v", err)
	}

	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()

	proto.RegisterCatalogServiceServer(server, dbService)

	reflection.Register(server)

	fmt.Println("Starting gRPC server on port :50051")
	if err := server.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
