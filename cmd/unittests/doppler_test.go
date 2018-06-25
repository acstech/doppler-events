package doppler_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Shopify/sarama"
	cb "github.com/acstech/doppler-events/internal/couchbase"
	"github.com/joho/godotenv"
)

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
	producer    sarama.AsyncProducer
)

func init() {
	//TODO check if development mode or not?????
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		panic(err)
	}
	//get config variables
	cbConn := os.Getenv("COUCHBASE_CONN")
	testCB := &cb.Couchbase{Doc: &cb.Doc{}}
	err = testCB.ConnectToCB(cbConn)
	if err != nil {
		//TODO try to manually create CB document, format, recheck error

	}

	config := sarama.NewConfig()
	brokers := []string{"localhost:9092"}
	producer, err = sarama.NewAsyncProducer(brokers, config)
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
