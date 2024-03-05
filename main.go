package main

import (
	"CatalogService/proto"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

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
