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

func Northeast() (float64, float64) {
	lat := 35 + rand.Float64()*(44-35)
	lng := 74 + rand.Float64()*(83-74)
	lng = lng - lng - lng
	return lat, lng
}

func Southeast() (float64, float64) {
	lat := 31 + rand.Float64()*(35-31)
	lng := 79 + rand.Float64()*(87-79)
	lng = lng - lng - lng
	return lat, lng
}

func Midwest() (float64, float64) {
	lat := 29 + rand.Float64()*(48-29)
	lng := 90 + rand.Float64()*(108-90)
	lng = lng - lng - lng
	return lat, lng
}

func West() (float64, float64) {
	lat := 33 + rand.Float64()*(48-33)
	lng := 108 + rand.Float64()*(121-108)
	lng = lng - lng - lng
	return lat, lng
}

func Random() (float64, float64) {
	lat := 26 + rand.Float64()*(48-26)
	lng := 80 + rand.Float64()*(118-80)
	lng = lng - lng - lng
	return lat, lng
}

func main() {
	//testing variables
	// numClients := 1 //the number of clients per test
	// numEvents := 1  //the number of events per client TODO make random number of events

	//data point variables
	var clientIDs = []string{"client0", "client1"}
	var eventIDs = []string{"physical check in", "mobile login", "rest"}
	//connect to server
	c, err := dial("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	//get true random
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var lng, lat float64
	for count := 0; count < 1500; count++ {
		time.Sleep(10 * time.Millisecond)
		clientID := clientIDs[r.Int31n(int32(len(clientIDs)))] //pick random client from clientIDs slice
		eventID := eventIDs[r.Int31n(int32(len(eventIDs)))]    //pick random event from eventIDs slice
		if count < 200 {
			lat, lng = Northeast()
		} else if count < 400 {
			if count%5 == 0 {
				lat, lng = Northeast()
			} else {
				lat, lng = Southeast()
			}
		} else if count < 600 {
			if count%5 == 0 {
				lat, lng = Southeast()
			} else if count%10 == 0 {
				lat, lng = Northeast()
			} else {
				lat, lng = Midwest()
			}
		} else if count < 1000 {
			if count%5 == 0 {
				lat, lng = Midwest()
			} else if count%10 == 0 {
				lat, lng = Southeast()
			} else if count%15 == 0 {
				lat, lng = Northeast()
			} else {
				lat, lng = West()
			}
		} else {
			lat, lng = Random()
		}
		resp, err := sendRequest(c, clientID, eventID, lat, lng)
		//lo any error
		if err != nil {
			log.Println(err)
			continue
		}
		//print server response
		log.Println(resp.Response, count)
	}
}

// func range(f1 float64, f2 float64) float64 {
//
// }

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
