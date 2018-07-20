package doppler_test

import (
	"testing"
)

var (
// bucket      string
// clientID    string
// badclientID string
// eventID     string
// dataName1   string
// dataName2   string
// dataSet     map[string]string
// cbconfig    string
// badcbconfig string
)

func init() {
	//TODO check if development mode or not?????

}

//TODO: Test the following

/*

ServeData
Init
DisplayData
ConnectToCB


*/

func TestInit(t *testing.T) {

	//TODO: Message structure
	//TODO: couchbase config
	//TODO: testsend data
	//TODO: find way to check error to determine if test passes or fails

	//if(err := Initialize != nil)
	//{
	// t.Errorf("someting wrong")
	// }
}

func TestInitFail(t *testing.T) {
	//TODO: some bad data
	//TODO: some invalid configs
	//TODO verify error response

	//t.Errorf("something right????")
}

func TestBatch(t *testing.T) {
	//  prod :=

}

func TestBatchFail(t *testing.T) {

}
