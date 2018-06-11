// Event API Client
//
//
//
//
//

package client

import (
	"log"
	//	"os"

	"time"

	pb "github.com/acstech/doppler-events/eventAPI"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "localhost:8080"
	// address     = "10.22.97.107:8080"
	defaultName = "default"
)

func DisplayData(clientID *string, eventID *string, dataSet *map[string]string) {

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewEventSenderClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SendEvent(ctx, &pb.EventObj{ClientID: *clientID, EventID: *eventID, DataSet: *dataSet})
	if err != nil {
		log.Fatalf("could not do anything: %v", err)
	}
	log.Printf(r.Response)
}
