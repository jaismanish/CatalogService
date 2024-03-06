package main

import (
	"CatalogService/proto"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

const (
	testDBHost     = "localhost"
	testDBPort     = "5432"
	testDBUser     = "postgres"
	testDBPassword = "Manish@2001"
	testDBName     = "CatalogTest"
)

var (
	testServer *grpc.Server
	testClient proto.CatalogServiceClient
	testDB     *sql.DB
)

func setupTest() {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		testDBHost, testDBPort, testDBUser, testDBPassword, testDBName))
	if err != nil {
		log.Fatal(err)
	}
	testDB = db

	err = createTables(db)
	if err != nil {
		log.Fatal(err)
	}

	catalogSrv := &catalogService{db: db}
	testServer = grpc.NewServer()
	proto.RegisterCatalogServiceServer(testServer, catalogSrv)

	go func() {
		listen, err := net.Listen("tcp", ":50052")
		if err != nil {
			log.Fatal(err)
		}
		if err := testServer.Serve(listen); err != nil {
			log.Fatal(err)
		}
	}()

	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	testClient = proto.NewCatalogServiceClient(conn)
}

func tearDownTest() {
	testServer.Stop()
	testDB.Close()
}

func TestMain(m *testing.M) {
	setupTest()
	code := m.Run()
	tearDownTest()
	os.Exit(code)
}

func TestAddRestaurant(t *testing.T) {
	req := &proto.AddRestaurantRequest{
		Name:     "Test Restaurant",
		Location: "Test Location",
	}

	res, err := testClient.AddRestaurant(context.Background(), req)
	assert.NoError(t, err)
	assert.True(t, res.Success)
}
