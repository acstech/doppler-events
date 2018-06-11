package main

import (
  "time"
  "log"
  "net"
  "golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "github.com/acstech/doppler-events/eventAPI"
	"google.golang.org/grpc/reflection"
)

const (
  // addr = "10.22.97.107:8080"
  addr = ":8080"
)

type server struct{}

func (s *server) SendEvent(ctx context.Context, in *pb.EventObj) (*pb.EventResp, error) {
  return &pb.EventResp{Response: "\nClient ID: " + in.ClientID + "\n" + "Event ID: " + in.EventID + "\nDate: " + time.Unix(int64(in.TimeSinceEpoch), 0).String() +"\nLongitude: " + in.KeyValues["lon"] + "\nLatitude: " + in.KeyValues["lat"]}, nil
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
