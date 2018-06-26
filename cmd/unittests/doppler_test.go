package doppler_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Shopify/sarama"
	cb "github.com/acstech/doppler-events/internal/couchbase"
	pb "github.com/acstech/doppler-events/rpc/eventAPI" //c meaning client call
	"github.com/golang/protobuf/ptypes"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type InfluxResponse struct {
	values []interface{}
}

type Responses struct {
	ValArray []InfluxResponse
}

var (
	bucket      string
	clientID    string
	badclientID string
	eventID     string
	dataName1   string
	dataName2   string
	cbconfig    string
	badcbconfig string
	cbConn      string
	producer    sarama.AsyncProducer
	client      pb.EventAPIClient
	testCB      *cb.Couchbase
)

func init() {
	//TODO check if development mode or not?????
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		panic(err)
	}
	//get config variables
	cbConn = os.Getenv("COUCHBASE_CONN")
	testCB = &cb.Couchbase{Doc: &cb.Doc{}}
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure()) //WithInsecure meaning no authentication required
	if err != nil {
		fmt.Errorf("Did not connect: %v", err)
	}

	client = pb.NewEventAPIClient(conn)

	clientID = "client0"
	eventID = "event0"
	//TODO: setup influx client
	// var fluxClient client.Client
	// q := fluxClient.NewQuery("select * from dopplerDataHistory", "doppler", "s")
	// if response, err := fluxClient.Query(q); err == nil && response.Error() == nil {
	// 	//fmt.Println(response.Results)
	//
	// 	for r := range response.Results {
	// 		var ifx InfluxResponse
	// 		ifx.values = response.Results[r].Series[r].Values[r]
	// 		res.ValArray = append(res.ValArray, ifx)
	// 	}
	// }
}

//TODO: Test the following

/*

ServeData
Init
DisplayData
ConnectToCB


*/
//
// func TestInit(t *testing.T) {
//
// 	//TODO: Message structure
// 	//TODO: couchbase config
// 	//TODO: testsend data
// 	//TODO: find way to check error to determine if test passes or fails
//
// 	//if(err := Initialize != nil)
// 	//{
// 	// t.Errorf("someting wrong")
// 	// }
// }
//
// func TestInitFail(t *testing.T) {
// 	//TODO: some bad data
// 	//TODO: some invalid configs
// 	//TODO verify error response
//
// 	//t.Errorf("something right????")
// }
//
// func TestBatch(t *testing.T) {
// 	//  prod :=
//
// }
//
// func TestBatchFail(t *testing.T) {
//
// }
//
// func TestInvalidData(t *testing.T) {
//
// }
func TestCBconnect(t *testing.T) {
	err := testCB.ConnectToCB(cbConn)
	if err != nil {
		//TODO try to manually create CB document, format, recheck error
		t.Errorf("success expected: this, received: %d", err)
	}
}

func TestCBconnect2(t *testing.T) {
	if !testCB.ClientExists(clientID) {
		fmt.Println("err")
	}
	//ensure that the eventID exists
	// testCB.EventEnsure(clientID, eventID)
	// if err != nil {
	// 	//TODO try to manually create CB document, format, recheck error
	// 	t.Errorf("success expected: this, received: %d", err)
	// }
}

func TestValidData(t *testing.T) {
	fmt.Println("testValid")
	dataSet := make(map[string]string, 2)
	dateTime := ptypes.TimestampNow()
	//dataSet["lat"] = 23.32
	//dataSet["lng"] = 32.23
	latitude := "32.23"
	longitude := "23.32"
	dataSet["lat"] = latitude
	dataSet["lng"] = longitude

	//req := &pb.DisplayRequest{ClientId: clientID, EventId: eventID, DateTime: dateTime, DataSet: dataSet}

	_, err := client.DisplayData(context.Background(), &pb.DisplayRequest{
		ClientId: clientID,
		EventId:  eventID,
		DateTime: dateTime,
		DataSet:  dataSet,
	})

	if err != nil {
		fmt.Println(err)
		t.Errorf("errored %d", err)
	}
}

func TestCleanup(t *testing.T) {

}

func TestCBevents(t *testing.T) {

}
