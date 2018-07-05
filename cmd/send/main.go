package main

import (
	"fmt"
	"net"
	"os"

	"github.com/acstech/doppler-events/internal/service"
	pb "github.com/acstech/doppler-events/rpc/eventAPI"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	// Load env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		// Have to load env file, so panic if you can't
		panic(err)
	}

	// Get config variables
	connString := os.Getenv("COUCHBASE_CONN")

	// Create an instance of event service
	eventService, err := service.Init()
	if err != nil {
		// Have to create a service instance, so panic if you can't
		panic(err)
	}

	// Initialize listener on server address
	lis, err := net.Listen("tcp", service.Address)
	if err != nil {
		fmt.Println("failed to create Kafka producer connection: ", err)
	}

	// Initialize grpc server
	grpcServer := grpc.NewServer()

	// Connect to couchbase
	err = eventService.CbConn.ConnectToCB(connString)
	if err != nil {
		fmt.Println("CB connection error: ", err)
	}

	// Register service to grpc
	pb.RegisterEventAPIServer(grpcServer, eventService)

	// Tells server to process incoming messages, checks if it failed to serve
	if err = grpcServer.Serve(lis); err != nil {
		fmt.Println("Failed to serve: ", err)
	}
}
