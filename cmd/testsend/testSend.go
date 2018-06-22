// Test User
//
//
//
//
//

package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	pb "github.com/acstech/doppler-events/rpc/eventAPI"
	"github.com/golang/protobuf/ptypes" //c meaning client call
	"google.golang.org/grpc"
)

func main() {
	//testing variables
	numClients := 1 //the number of clients per test
	numEvents := 1  //the number of events per client TODO make random number of events

	//data point variables
	var clientID = "client0"
	var eventID = ""
	var dataName1 = "lat"
	var dataName2 = "lng"
	var dataSet = make(map[string]string)

	c, err := dial("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	// loop to create test data (i is number of data points)
	// first loop creates number of clients
	// TODO fix this function
	for clientNum := 0; clientNum < numClients; clientNum++ {
		//second loop creates number of events
		for eventNum := 0; eventNum < numEvents; eventNum++ {
			eventID = fmt.Sprint("event", eventNum)
			//third loop iterates through all latitudes and gets current time
			for lat := -85; lat <= 85; lat++ {
				dataSet[dataName1] = strconv.Itoa(lat)
				//fourth loop iterates through all longitudes and calls API
				for lng := -175; lng <= 175; lng++ {
					//time.Sleep(250 * time.Millisecond)

					dataSet[dataName2] = strconv.Itoa(lng)
					//fmt.Println(clientID, " ", eventID, " ", dateTime, " ", dataSet)
					dateTime := ptypes.TimestampNow() //get current time
					resp, err := c.DisplayData(context.Background(), &pb.DisplayRequest{
						ClientId: clientID,
						EventId:  eventID,
						DateTime: dateTime,
						DataSet:  dataSet,
					})
					if err != nil {
						log.Println(err)
						continue
					}
					log.Println(resp.Response)
				}
			}
		}
	}
}

func dial(addr string) (pb.EventAPIClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure()) //WithInsecure meaning no authentication required
	if err != nil {
		return nil, fmt.Errorf("Did not connect: %v", err)
	}
	client := pb.NewEventAPIClient(conn)
	return client, nil
}
