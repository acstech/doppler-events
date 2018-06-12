// Test User
//
//
//
//
//

package main

import (
	c "github.com/acstech/doppler-events/eventAPI_client" //c meaning client call
)

func main() {

	//sample data
	clientID := "Bob"
	eventID := "Sign In"

	dataName1 := "lon"
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
	c.DisplayData(&clientID, &eventID, &dataSet)

}
