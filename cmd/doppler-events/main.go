package main

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"

	cb "github.com/acstech/doppler-events/internal/couchbase"
	"github.com/acstech/doppler-events/internal/service"
	pb "github.com/acstech/doppler-events/rpc/eventAPI"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
)

func main() {
	// Create signals and done channel to handle graceful shutdown of server
	signals := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	// Wait for any signals to happen, such as interrupts
	go func() {
		sig := <-signals
		fmt.Println(sig)
		done <- true
	}()

	// Get config variables
	connString := os.Getenv("COUCHBASE_CONN")
	address := os.Getenv("API_ADDRESS")
	kafkaConn, kafkaTopic, err := kafkaParse(os.Getenv("KAFKA_CONN"))
	if err != nil {
		panic(err)
	}

	// Create asynchronous Kafka producer
	producer, err := newProducer(kafkaConn)
	if err != nil {
		panic(err)
	}

	// Create an empty couchbase connection instance
	cbConn := &cb.Couchbase{}

	// Connect to couchbase
	err = cbConn.ConnectToCB(connString)
	if err != nil {
		fmt.Println("CB connection error: ", err)
	}

	// Create an instance of event service
	eventService, err := service.NewService(producer, cbConn, kafkaTopic)
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

	// Graceful shutdown of producer, couchbase connection, and gRPC server
	defer func() {
		if errt := producer.Close(); errt != nil {
			err = fmt.Errorf("error closing producer: %v", errt)
		}
		if err = cbConn.Bucket.Close(); err != nil {
			err = fmt.Errorf("error closing couchbase connection: %v", err)
		}
		grpcServer.GracefulStop()
	}()

	// Register service to grpc
	pb.RegisterEventAPIServer(grpcServer, eventService)

	// Tells server to process incoming messages, checks if it failed to serve
	if err = grpcServer.Serve(lis); err != nil {
		fmt.Println("Failed to serve: ", err)
	}

	// Block until interrupt detected
	<-done
	fmt.Println("Exiting")
}

// Parses env variable to return kafka address and topic
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
