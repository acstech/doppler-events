// Event API Server
//

package service

import (
	"fmt"
	"net"

	"github.com/Shopify/sarama"
	cb "github.com/acstech/doppler-events/internal/couchbase"
	pb "github.com/acstech/doppler-events/rpc/eventAPI"
	"google.golang.org/grpc"
)

// Initialize addressess of where server listens with format IP:Port
const (
	address = ":8080"
)

// Struct to hold parameters for server
type Service struct {
	errorResponse ErrorRes
	producer      sarama.AsyncProducer
	cbConn        *cb.Couchbase
}

// Error struct
type ErrorRes struct {
	err    error
	errMes []string
}

//
type ServiceConfig struct {
	cbCon string
}

// Create an instance of Service
func New(e ErrorRes, p sarama.AsyncProducer, c *cb.Couchbase) *Service {
	return &Service{
		errorResponse: e,
		producer:      p,
		cbConn:        c,
	}
}

// New returns an error that formats as the given text.
// func New(text string) error {
// 	return &errorString{text}
// }

// errorString is a trivial implementation of error.
// type errorString struct {
// 	s string
// }

// func (e *errorString) Error() string {
// 	return e.s
// }

// Init sets up the backend server so that clients can send data to kafka
// cbCon is the connection string for couchbase
// returns an error if any occur while creating a kafka producer, a couchbase connection, sending data,
// or closing the kafka producer
func Init(cbCon, kafkaCon, kafkaTopic, address string) error {
	//initialize listener on server address
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	//initialize server
	grpcServer := grpc.NewServer()
	// Create a new producer
	prod, err := newProducer()
	if err != nil {
		return fmt.Errorf("failed to create Kafka producer connection: %v", err)
	}

	defer func() {
		if errt := prod.Close(); errt != nil {
			// Should not reach here
			err = fmt.Errorf("error closing producer: %v", errt)
		}
	}()
	// Instance of server
	serve2 := server{
		theProd:    prod,
		cbConn:     &cb.Couchbase{},
		kafkaTopic: kafkaTopic,
	}
	err = serve2.cbConn.ConnectToCB(cbCon)
	if err != nil {
		return fmt.Errorf("CB connection error: %v", err)
	}
	//register server to grpc
	pb.RegisterEventAPIServer(grpcServer, &serve2)
	// Register reflection service on gRPC server for back and forth communication. Kept for future use if necessary
	// reflection.Register(s)

	//tells the server to process the incoming messages, checks if failed to serve
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return err
}
