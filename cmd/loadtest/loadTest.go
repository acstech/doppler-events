<<<<<<< HEAD
<<<<<<< HEAD
//This test file is for a fast, infinite load test that does not clean up influx data after itself.
//The purpose of this file is speed

=======
>>>>>>> changed
=======
//This test file is for a fast, infinite load test that does not clean up influx data after itself.
//The purpose of this file is speed

>>>>>>> b1a7f4bbc454ed2a3595d1e4448f49841644cd6e
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	pb "github.com/acstech/doppler-events/rpc/eventAPI" //c meaning client call
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
)

var (
	clientIDs []string
	eventIDs  []string
	c         pb.EventAPIClient
)

func main() {
	//data point variables
	clientIDs = []string{"client0", "client1"} //In order for test to work, couchbase must contain all 3 clients
	eventIDs = []string{"physical check in", "mobile login", "rest"}
	var err error
	//connect to server
	c, err = dial("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("starting load test...")
	LoadTest()

}

//LoadTest sends infinite random points to the API
func LoadTest() {
	//get true random
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		//time.Sleep(3 * time.Millisecond)
		clientID := clientIDs[r.Int31n(int32(len(clientIDs)))]   //pick random client from clientIDs slice
		eventID := eventIDs[r.Int31n(int32(len(eventIDs)))]      //pick random event from eventIDs slice
		lat := (r.Float64() - .5) * 180                          //get random lat
		lng := (r.Float64() - .5) * 360                          //get random lng
		resp, err := sendRequest(c, clientID, eventID, lat, lng) //call function that prepares data to send to server
		//lo any error
		if err != nil {
			log.Println(err)
			continue
		}
		//print server response
		if resp != nil {
			log.Println(resp.Response)
		}
	}
}

//takes client data, sends it to over connection
func sendRequest(c pb.EventAPIClient, clientID, eventID string, lat, lng float64) (*pb.DisplayResponse, error) {
	//create map of data
	dataSet := make(map[string]string, 2)
	dataSet["lat"] = strconv.FormatFloat(lat, 'g', -1, 64)
	dataSet["lng"] = strconv.FormatFloat(lng, 'g', -1, 64)
	//get current time
	dateTime := ptypes.TimestampNow()
	//send data to server, returns response and error

	resp, err := c.DisplayData(context.Background(), &pb.DisplayRequest{
		ClientId: clientID,
		EventId:  eventID,
		DateTime: dateTime,
		DataSet:  dataSet,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

//get grpc connection client
func dial(addr string) (pb.EventAPIClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure()) //WithInsecure meaning no authentication required
	if err != nil {
		return nil, fmt.Errorf("Did not connect: %v", err)
	}
	client := pb.NewEventAPIClient(conn)
	return client, nil
}
