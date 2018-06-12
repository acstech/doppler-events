// Event API Client
//
//
//
//
//
package main

import (
	"log"
	"time"
	//	"os"
	pb "github.com/acstech/doppler-events/eventAPI"
	"github.com/golang/protobuf/ptypes"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "localhost:8080"
	// address     = "address to server"
)

func main() {
	//sample data
	clientID := "Bob"
	eventID := "Sign In"
	dataName1 := "lon"
	dataValue1 := "665.67676"
	dataName2 := "lat"
	dataValue2 := "87687"
	//create map of data
	dataSet := make(map[string]string)
	dataSet[dataName1] = dataValue1
	dataSet[dataName2] = dataValue2
	//call to API, requires clientID, eventID, and a dataSet
	DisplayData(&clientID, &eventID, &dataSet)
}

func DisplayData(clientID *string, eventID *string, dataSet *map[string]string) {

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
	r, err := c.SendEvent(ctx, &pb.EventObj{ClientID: *clientID, EventID: *eventID, TimeSinceEpoch: ptypes.TimestampNow(), DataSet: *dataSet})
	if err != nil {
		log.Fatalf("could not do anything: %v", err)
	}
	//print out server reponse to console
	log.Printf(r.Response)
}
