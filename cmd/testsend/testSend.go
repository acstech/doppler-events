// Test User
//
//
//
//
//

package main

import (
	"fmt"
	"strconv"
	"time"

	c "github.com/acstech/doppler-events/cmd/grpcTEST/EventAPIClient"
	"github.com/golang/protobuf/ptypes" //c meaning client call
)

func main() {
	//testing variables
	numClients := 1 //the number of clients per test
	numEvents := 1  //the number of events per client TODO make random number of events

	//data point variables
	var clientID = "test2"
	var eventID = "test"
	var dataName1 = "lat"
	var dataName2 = "lng"
	var dataSet = make(map[string]string)

	// loop to create test data (i is number of data points)
	// first loop creates number of clients
	// TODO fix this function
	for clientNum := 0; clientNum < numClients; clientNum++ {
		clientID = fmt.Sprint("client", clientNum)
		//second loop creates number of events
		for eventNum := 0; eventNum < numEvents; eventNum++ {
			eventID = fmt.Sprint("event", eventNum)
			//third loop iterates through all latitudes and gets current time
			for lat := -85; lat <= 85; lat++ {
				dataSet[dataName1] = strconv.Itoa(lat)
				//fourth loop iterates through all longitudes and calls API
				for lng := -175; lng <= 175; lng++ {
					time.Sleep(500 * time.Millisecond)
					dataSet[dataName2] = strconv.Itoa(lng)
					//fmt.Println(clientID, " ", eventID, " ", dateTime, " ", dataSet)
					dateTime := ptypes.TimestampNow() //get current time
					c.DisplayData(clientID, eventID, dateTime, dataSet)
				}
			}
		}
	}
}
