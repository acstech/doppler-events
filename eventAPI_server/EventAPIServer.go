package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/acstech/doppler-events/eventAPI"
	"github.com/golang/protobuf/ptypes"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	addr = "10.22.97.107:8080"
	// addr = ":8080"
)

type server struct{}

func (s *server) SendEvent(ctx context.Context, in *pb.EventObj) (*pb.EventResp, error) {
	ts, err := ptypes.Timestamp(in.TimeSinceEpoch)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Client ID: %s\nEvent ID: %s\nLongitude: %s\nLatitude: %s\n", in.ClientID, in.EventID, in.KeyValues["lon"], in.KeyValues["lat"])
	fmt.Println("Date: ", ts.String())
	return &pb.EventResp{Response: "\nClient ID: " + in.ClientID + "\nEvent ID: " + in.EventID + "\nDate: " + ts.String() + "\nLongitude: " + in.KeyValues["lon"] + "\nLatitude: " + in.KeyValues["lat"]}, nil
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
