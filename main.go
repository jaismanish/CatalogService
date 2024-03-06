package main

import (
	"CatalogService/proto"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

const (
	dbHost     = "localhost"
	dbPort     = "5432"
	dbUser     = "postgres"
	dbPassword = "Manish@2001"
	dbName     = "Catalog"
)

type catalogService struct {
	db *sql.DB
	proto.UnimplementedCatalogServiceServer
}

func NewCatalogService() (*catalogService, error) {
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

	return &catalogService{
		db: db,
	}, nil
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

func (s *catalogService) AddRestaurant(ctx context.Context, req *proto.AddRestaurantRequest) (*proto.AddRestaurantResponse, error) {
	_, err := s.db.Exec("INSERT INTO restaurants (name, location) VALUES ($1, $2)", req.Name, req.Location)
	if err != nil {
		return nil, fmt.Errorf("failed to add restaurant: %v", err)
	}

	return &proto.AddRestaurantResponse{Success: true}, nil
}

func (s *catalogService) AddMenuItem(ctx context.Context, req *proto.AddMenuItemRequest) (*proto.AddMenuItemResponse, error) {
	return &proto.AddMenuItemResponse{Success: true}, nil
}

func (s *catalogService) GetRestaurants(ctx context.Context, req *proto.GetRestaurantsRequest) (*proto.GetRestaurantsResponse, error) {
	return &proto.GetRestaurantsResponse{}, nil
}

func (s *catalogService) GetMenuItems(ctx context.Context, req *proto.GetMenuItemsRequest) (*proto.GetMenuItemsResponse, error) {
	return &proto.GetMenuItemsResponse{}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()

	if err != nil {
		log.Fatalf("Failed to create catalog service: %v", err)
	}

	fmt.Println("Starting gRPC server on port :50051")
	if err := server.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
