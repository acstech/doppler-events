// Test User
//
//
//
//
//

package main

import (
	c "github.com/acstech/doppler-events/eventAPI_client"
)

func main() {

	//sample data
	clientID := "Bob"
	eventID := "Sign In"

	dataName1 := "lon"
	dataValue1 := "665.67676"
	dataName2 := "lat"
	dataValue2 := "87687"

	dataSet := make(map[string]string)
	dataSet[dataName1] = dataValue1
	dataSet[dataName2] = dataValue2

	//call to API
	c.DisplayData(&clientID, &eventID, &dataSet)

}
