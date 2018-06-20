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

	c "github.com/acstech/doppler-events/cmd/grpcTEST/EventAPIClient" //c meaning client call
	ptypes "github.com/golang/protobuf/ptypes"
)

func main() {

	//testing variables
	numClients := 1 //the number of clients per test
	numEvents := 1  //the number of events per client TODO make random number of events

	//data point variables
	var clientID string
	var eventID string
	var dataName1 = "lat"
	var dataName2 = "lng"
	var dataSet = make(map[string]string)

	//loop to create test data (i is number of data points)
	//first loop creates number of clients
	for clientNum := 0; clientNum < numClients; clientNum++ {
		clientID = fmt.Sprint("client", clientNum)
		//second loop creates number of events
		for eventNum := 0; eventNum < numEvents; eventNum++ {
			eventID = fmt.Sprint("event", eventNum)
			//third loop iterates through all latitudes and gets current time
			for lat := -90; lat <= 90; lat++ {
				dateTime := ptypes.TimestampNow() //get current time
				dataSet[dataName1] = strconv.Itoa(lat)
				//fourth loop iterates through all longitudes and calls API
				for lng := -180; lng <= 180; lng++ {
					dataSet[dataName2] = strconv.Itoa(lng)
					//fmt.Println(clientID, " ", eventID, " ", dateTime, " ", dataSet)
					c.DisplayData(clientID, eventID, dateTime, dataSet)
				}
			}
		}
	}

}
