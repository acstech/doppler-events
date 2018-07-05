// Event API Server

package service

import (
	"fmt"
	//"net"

	"github.com/Shopify/sarama"
	cb "github.com/acstech/doppler-events/internal/couchbase"
	//pb "github.com/acstech/doppler-events/rpc/eventAPI"
	//"google.golang.org/grpc"
)

// Initialize addressess of where server listens with format IP:Port
const (
	Address = ":8080"
)

// Struct to hold parameters for server
type Service struct {
	errorResponse ErrorRes
	producer      sarama.AsyncProducer
	CbConn        *cb.Couchbase
}

// Error struct for verifying constraints
type ErrorRes struct {
	err    error
	errMes []string
}

// Init sets up the backend server so that clients can send data to kafka
// cbCon is the connection string for couchbase
// returns an error if any occur while creating a kafka producer, a couchbase connection, sending data,
// or closing the kafka producer
func Init() (*Service, error) {
	e := ErrorRes{}

	// Initialize listener on server address
	// lis, err := net.Listen("tcp", address)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to listen: %v", err)
	// }

	// Initialize grpc server
	//grpcServer := grpc.NewServer()

	// Create a new Kafka producer
	prod, err := newProducer()
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer connection: %v", err)
	}

	// Close after everything else has finished executing
	defer func() {
		if errt := prod.Close(); errt != nil {
			// Should not reach here
			err = fmt.Errorf("error closing producer: %v", errt)
		}
	}()

	// err = serve2.cbConn.ConnectToCB(cbCon) // Call in main func
	// if err != nil {
	// 	return fmt.Errorf("CB connection error: %v", err)
	// }

	//register server to grpc
	//pb.RegisterEventAPIServer(grpcServer, &serve2)
	// Register reflection service on gRPC server for back and forth communication. Kept for future use if necessary
	// reflection.Register(s)

	//tells the server to process the incoming messages, checks if failed to serve
	// if err := grpcServer.Serve(lis); err != nil {
	// 	return nil, fmt.Errorf("failed to serve: %v", err)
	// }

	return &Service{
		errorResponse: e,
		producer:      prod,
		CbConn:        &cb.Couchbase{},
	}, nil
}
