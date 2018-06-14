// Event API Client
// Files Needed: EventAPIClient.go, EventAPIServer.go, eventAPI.pb.go
// About: this program takes in a clientID *string, eventID *string, and a *map[dataName string]dataValue string
// and enables the data to be displayed on the heatmap found at _____ ********
//
// Specifically this program sends data to Kafka so that it can be processed before displaying
//
// How to test:
// start server with "go run eventAPI_server"
// start client with "go run testUser.go"
//
// How to use:
// Import this package via "github.com/acstech/doppler-events/eventAPI_client"
// Initialize clientID and eventID
// Initialize a map with the dataName as the key and the dataValue as the value
// Call c.DisplayData(clientID *string, eventID *string, dataSet *map[dataName string]dataValue string with the required arguments
//

package client

import (
	"log"
	"time"
	//	"os"
	pb "github.com/acstech/doppler-events/eventAPI"
	"github.com/golang/protobuf/ptypes"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

//intialize address of Kafka connect
const (
	address = "localhost:8080" //TODO: <- this isn't kafka's address?
	// address = "kafka IP:Port"
)

//DisplayData takes in the required arguments and sends the arguments to the server
func DisplayData(clientID *string, eventID *string, dataSet *map[string]string) {

	//intialize connection with server at address
	conn, err := grpc.Dial(address, grpc.WithInsecure()) //WithInsecure meaning no authentication required
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}

	defer conn.Close()                 //wait until DisplayData function returns then closes connection
	c := pb.NewEventSenderClient(conn) //initializes new conneciton to server, calls to pb.go

	//set up connection context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//send event data to server, save response
	r, err := c.SendEvent(ctx, &pb.EventObj{ClientId: *clientID, EventId: *eventID, timeSinceEpoch: ptypes.TimestampNow(), DataSet: *dataSet})
	if err != nil {
		log.Fatalf("could not do anything: %v", err)
	}
	//print out server reponse to console
	log.Printf(r.Response)
}
