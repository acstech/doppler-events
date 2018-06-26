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
	"math/rand"
	"strconv"
	"time"

	pb "github.com/acstech/doppler-events/rpc/eventAPI" //c meaning client call
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
)

func main() {
	//testing variables
	// numClients := 1 //the number of clients per test
	// numEvents := 1  //the number of events per client TODO make random number of events

	//data point variables
	var clientIDs = []string{"client0", "nav1"}
	var eventIDs = []string{"event0", "event1"}
	//connect to server
	c, err := dial("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	//get true random
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	count := 0
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
		count++
		log.Println(resp.Response, count)
	}

	// loop to create test data (i is number of data points)
	// first loop creates number of clients
	// for clientNum := 0; clientNum < numClients; clientNum++ {
	// 	//second loop creates number of events
	// 	for eventNum := 0; eventNum < numEvents; eventNum++ {
	// 		eventID = fmt.Sprint("event", eventNum)
	// 		//third loop iterates through all latitudes and gets current time
	// 		for lat := 0; lat <= 85; lat++ {
	// 			dataSet[dataName1] = strconv.Itoa(lat)
	// 			//fourth loop iterates through all longitudes and calls API
	// 			for lng := 0; lng <= 175; lng++ {
	// 				dataSet[dataName2] = strconv.Itoa(lng)
	// 				//fmt.Println(clientID, " ", eventID, " ", dateTime, " ", dataSet)
	// 				dateTime := ptypes.TimestampNow() //get current time
	// 				resp, err := c.DisplayData(context.Background(), &pb.DisplayRequest{
	// 					ClientId: clientID,
	// 					EventId:  eventID,
	// 					DateTime: dateTime,
	// 					DataSet:  dataSet,
	// 				})
	// 				if err != nil {
	// 					log.Println(err)
	// 					continue
	// 				}
	// 				log.Println(resp.Response)
	// 			}
	// 		}
	// 	}
	// }
}

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

func dial(addr string) (pb.EventAPIClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure()) //WithInsecure meaning no authentication required
	if err != nil {
		return nil, fmt.Errorf("Did not connect: %v", err)
	}
	client := pb.NewEventAPIClient(conn)
	return client, nil
}
