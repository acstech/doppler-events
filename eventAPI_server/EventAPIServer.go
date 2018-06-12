// Event API Server
// Files Needed: EventAPIClient.go, EventAPIServer.go, eventAPI.pb.go
// About: this program receives a message from EventAPIClient.go, formats the data to JSON, and sends the data to Kafka Connect
//

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/acstech/doppler-events/eventAPI"
	"github.com/golang/protobuf/ptypes"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

//intialize addressess of where server listens with format IP:Port
const (
	// address = "10.22.97.107:8080"
	address = ":8080"
)

//struct to hold parameters for server
type server struct{}

// SendEvent is the function that EventAPIClient.go calls in order to send data to the server
// the data is then processed, formatted to JSON, and send to Kafka Connect
func (s *server) SendEvent(ctx context.Context, in *pb.EventObj) (*pb.EventResp, error) {
	//get the current time
	ts, err := ptypes.Timestamp(in.TimeSinceEpoch)
	if err != nil {
		return nil, err
	}

	//printing ClientID and EventID to server console for testing
	fmt.Printf("ClientID: " + in.ClientID + "\nEventID: " + in.EventID + "\nDate: " + ts.String() + "\n")

	//printing DataSet to server console for testing
	for key, value := range in.DataSet {
		fmt.Println("DataType: ", key, "DataValue:", value)
	}
	fmt.Println()

	//format to JSON
	bytes, err := json.Marshal(in) //Marshal returns the ascii presentation of the data
	if err != nil {
		fmt.Println("Format to JSON Error")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	//print JSON format of dataSet to console for testing
	fmt.Println(string(bytes))

	//NEEDED method to send to Kafka

	//return response to client
	return &pb.EventResp{Response: "Success! Open heatmap at ____ to see results"}, nil
}

func main() {
	//initialize listener on server address
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	//initialize server
	s := grpc.NewServer()
	//register server to grpc
	pb.RegisterEventSenderServer(s, &server{})
	// Register reflection service on gRPC server for back and forth communication. Kept for future use if necessary
	// reflection.Register(s)

	//tells the server to process the incoming messages, checks if failed to serve
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
