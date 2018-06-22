package doppler_test

import "testing"

var (
	bucket      string
	clientID    string
	badclientID string
	eventID     string
	dataName1   string
	dataName2   string
	dataSet     map[string]string
	cbconfig    string
	badcbconfig string
)

func init() {

}

//TODO: Test the following

/*

ServeData
Init
DisplayData
ConnectToCB


*/

func Test1(t *testing.T) {
	//TODO: Message structure
	//TODO: couchbase config
	//TODO: testsend data
	//TODO: find way to check error to determine if test passes or fails

	//if(functionresult==bad)
	//{
	// t.Errorf("someting wrong")
	// }
}

func TestFail1(t *testing.T) {
	//TODO: some bad data
	//TODO: some invalid configs
	//TODO verify error response
}

func Test2(t *testing.T) {

}
