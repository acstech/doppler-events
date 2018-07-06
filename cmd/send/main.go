package main

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"

	"github.com/acstech/doppler-events/internal/service"
	pb "github.com/acstech/doppler-events/rpc/eventAPI"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
)

func main() {
	// Get config variables
	connString := os.Getenv("COUCHBASE_CONN")
	address := os.Getenv("API_ADDRESS")
	fmt.Println("HERE")
	kafkaConn, kafkaTopic, err := kafkaParse(os.Getenv("KAFKA_CONN"))
	if err != nil {
		panic(err)
	}

	// Create an instance of event service
	eventService, err := service.Init(kafkaConn, kafkaTopic, address)
	if err != nil {
		// Have to create a service instance, so panic if you can't
		panic(err)
	}

	// Initialize listener on server address
	lis, err := net.Listen("tcp", address)
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

func kafkaParse(conn string) (string, string, error) {
	u, err := url.Parse(conn)
	if err != nil {
		return "", "", err
	}
	fmt.Println("HERE: ", u.Host)
	if u.Host == "" {
		return "", "", errors.New("Kafka address is not specified, verify that your environment varaibles are correct")
	}
	address := u.Host
	// make sure that the topic is specified
	if u.Path == "" || u.Path == "/" {
		return "", "", errors.New("Kafka topic is not specified, verify that your environment varaibles are correct")
	}
	topic := u.Path[1:]
	return address, topic, nil
}
