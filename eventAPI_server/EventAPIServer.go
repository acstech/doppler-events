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
	"strconv"

	pb "github.com/acstech/doppler-events/eventAPI"
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

	//converting protobuf timestap to seconds in unix time to a string
	ts := strconv.FormatInt(in.DateTime.GetSeconds(), 10)

	//convert EventObj to map in order to flatten (needed to flatten for influxDB)
	//intialize flatJSONMap as placeholder for marshal
	flatJSONMap := make(map[string]string)
	//will always have clientID, eventID, dateTime
	flatJSONMap["clientID"] = in.ClientId
	flatJSONMap["eventID"] = in.EventId
	flatJSONMap["dateTime"] = ts
	//loop across dataSet map and add key and value to flatJSON
	for key, value := range in.DataSet {
		flatJSONMap[key] = value
	}

	//format to JSON
	JSONbytes, err := json.Marshal(flatJSONMap) //Marshal returns the ascii presentation of the data
	if err != nil {
		fmt.Println("Format to JSON Error")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	//print JSON to server console for testing
	fmt.Println(string(JSONbytes))

	//NEEDED method to send to Kafka
	// sendToKafka(string(JSONbytes))

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
