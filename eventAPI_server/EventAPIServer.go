// Event API Server
//
//
//
//

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	// addr = "10.22.97.107:8080"
	addr = ":8080"
)

type server struct{}

func (s *server) SendEvent(ctx context.Context, in *pb.EventObj) (*pb.EventResp, error) {
	//printing ClientID and EventID to server console
	ts, err := ptypes.Timestamp(in.TimeSinceEpoch)
	if err != nil {
		return nil, err
	}
	fmt.Printf("ClientID: " + in.ClientID + "\nEventID: " + in.EventID + "\nDate: " + ts.String() "\n")

	//printing DataSet to server console
	for key, value := range in.DataSet {
		fmt.Println("DataType: ", key, "DataValue:", value)
	}
	fmt.Println()

	//method to format to JSON
	bytes, err := json.Marshal(in) //Marshal returns the ascii presentation of the data
	if err != nil {
		fmt.Println("Format to JSON Error")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	//print JSON format of
	fmt.Println(string(bytes))

	//method to send to Kafka

	//return response to client
	return &pb.EventResp{Response: "Success! Open heatmap at ____ to see results"}, nil
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
