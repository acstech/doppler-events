// Event API Server
//
//
//
//

package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/acstech/doppler-events/eventAPI"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	// addr = "10.22.97.107:8080"
	addr = ":8080"
)

type server struct{}

func (s *server) SendEvent(ctx context.Context, in *pb.EventObj) (*pb.EventResp, error) {
	fmt.Printf("ClientID: " + in.ClientID + "\nEventID: " + in.EventID + "\n")
	for key, value := range in.DataSet {
		fmt.Println("DataType: ", key, "DataValue:", value)
	}
	fmt.Println()
	//method to format to JSON
	//method to send to Kafka
	return &pb.EventResp{Response: "Client ID: " + in.ClientID}, nil
}

func main() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterEventSenderServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
