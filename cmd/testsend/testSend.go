// Test User
//
//
//
//
//

package main

import (
	ptype "github.com/golang/protobuf/ptypes"

	c "github.com/acstech/doppler-events/cmd/grpcTEST/EventAPIClient" //c meaning client call
)

func main() {

	ts := ptype.TimestampNow()
	//sample data
	clientID := "test2"
	eventID := "Sign In"
	dateTime := ts

	dataName1 := "lng"
	dataValue1 := "665.67676"
	dataName2 := "lat"
	dataValue2 := "87687"
	// dataName3 := "gender"
	// dataValue3 := "male"

	//create map of data
	dataSet := make(map[string]string)
	dataSet[dataName1] = dataValue1
	dataSet[dataName2] = dataValue2
	// dataSet[dataName3] = dataValue3

	//call to API, requires clientID, eventID, and a dataSet
	// service.Init()
	c.DisplayData(clientID, eventID, dateTime, dataSet)

}
