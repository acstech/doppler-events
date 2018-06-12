// Event API Client
//
//
//
//
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

const (
	address = "localhost:8080"
	// address     = "address to server"
)

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
