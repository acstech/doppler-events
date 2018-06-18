// Event API Client
//

package EventAPIClient

import (
	"log"
	"time"
	//	"os"
	pb "github.com/acstech/doppler-eventsTest/rpc/eventAPI"
	ptype "github.com/golang/protobuf/ptypes/timestamp"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

//intialize address of Kafka connect
const (
	address = "localhost:8080" //TODO: <- this isn't kafka's address?
	// address = "kafka IP:Port"
)

//DisplayData takes in the required arguments and sends the arguments to the server
func DisplayData(clientID string, eventID string, dateTime *ptype.Timestamp, dataSet map[string]string) {

	//intialize connection with server at address
	conn, err := grpc.Dial(address, grpc.WithInsecure()) //WithInsecure meaning no authentication required
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}

	defer conn.Close()              //wait until DisplayData function returns then closes connection
	c := pb.NewEventAPIClient(conn) //initializes new conneciton to server, calls to pb.go

	//set up connection context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//send event data to server, save response
	r, err := c.DisplayData(ctx, &pb.DisplayRequest{ClientId: clientID, EventId: eventID, DateTime: dateTime, DataSet: dataSet})
	if err != nil {
		log.Fatalf("could not do anything: %v", err)
	}
	//print out server reponse to console
	log.Printf(r.Response)
}
