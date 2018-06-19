// Event API Server
//

package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/Shopify/sarama"
	cb "github.com/acstech/doppler-events/internal/couchbase"
	pb "github.com/acstech/doppler-events/rpc/eventAPI"
	ptype "github.com/golang/protobuf/ptypes"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var cbHost, cbBucket, cbUsername, cbPassword string

type NewDisplayData struct {
	// DateTime int64  `json:"DateTime"`
	// ClientID string `json:"clientID"`
	// EventID  string `json:"eventID"`
	// Lng      string `json:"lng"`
	// Lat      string `json:"lat"`
}

//intialize addressess of where server listens with format IP:Port
const (
	address = ":8080"
)

//struct to hold parameters for server
type server struct {
	theProd sarama.AsyncProducer
}

//newProducer configures an asynchronous kafka producer client, returns it
func newProducer() (sarama.AsyncProducer, error) {
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	// Setup configuration
	config := sarama.NewConfig()
	config.ClientID = "1"
	// The level of acknowledgement reliability needed from the broker.
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	brokers := []string{"localhost:9092"}
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		// Should not reach here
		return nil, err
	}
	return producer, nil

}

//sendToQueue takes byte array, passes it to producer and writes to kafka instance
func (prod *server) sendToQueue(JSONob []byte) {

	var enqueued, errors int

	msg := &sarama.ProducerMessage{
		Topic: "influx-topic",
		Value: sarama.ByteEncoder(JSONob),
	}

	select {
	case prod.theProd.Input() <- msg:
		enqueued++
		fmt.Println("Produce message")
	case err := <-prod.theProd.Errors():
		errors++
		fmt.Println("Failed to produce message:", err)
	}

	log.Printf("Enqueued: %d; errors: %d\n", enqueued, errors)
}

// DisplayData is the function that EventAPIClient.go calls in order to send data to the server
// the data is then processed, formatted to JSON, and send to Kafka
func (s *server) DisplayData(ctx context.Context, in *pb.DisplayRequest) (*pb.DisplayResponse, error) {

	//converting protobuf timestap to to a string in format yyyy-MM-DDTHH:mm:ss.SSSZ
	ts := ptype.TimestampString(in.DateTime)

	//convert DisplayRequest to map in order to flatten (needed to flatten for influxDB)
	//intialize flatJSONMap as placeholder for marshal
	flatJSONMap := make(map[string]string)
	//check to make sure that the ClientID exists
	cbConn := &cb.Couchbase{Doc: &cb.Doc{}}
	cbConn.ConnectToCB(fmt.Sprintf("couchbase://%s:%s@%s/%s", cbUsername, cbPassword, cbHost, cbBucket))
	if !cbConn.ClientExists(in.ClientId) {
		return &pb.DisplayResponse{Response: "The ClientID is not valid."}, nil
	}
	//ensure that the eventID exists
	cbConn.EventEnsure(in.ClientId, in.EventId)
	fmt.Println("Client exists and the eventID has been ensured.")
	//will always have clientID, eventID, dateTime
	flatJSONMap["clientID"] = in.ClientId
	flatJSONMap["eventID"] = in.EventId
	flatJSONMap["dateTime"] = ts
	//loop across dataSet map and add key and value to flatJSON
	for key, value := range in.DataSet {
		flatJSONMap[key] = value
	}

	//format to JSON
	JSONbytes, err := json.Marshal(flatJSONMap) //Marshal returns the ascii presentation of the data
	if err != nil {
		fmt.Println("Format to JSON Error")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	//NEEDED method to send to Kafka
	s.sendToQueue(JSONbytes)

	//return response to client
	return &pb.DisplayResponse{Response: "Success! Open heatmap at ____ to see results"}, nil
}

func Init(cbH, cbBuck, cbUser, cbPass string) {
	//store config variables
	cbHost = cbH
	cbBucket = cbBuck
	cbUsername = cbUser
	cbPassword = cbPass
	//initialize listener on server address
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	//initialize server
	s := grpc.NewServer()
	prod, err := newProducer()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := prod.Close(); err != nil {
			// Should not reach here
			panic(err)
		}
	}()

	serve2 := server{
		theProd: prod,
	}
	//register server to grpc
	pb.RegisterEventAPIServer(s, &serve2)
	// Register reflection service on gRPC server for back and forth communication. Kept for future use if necessary
	// reflection.Register(s)

	//tells the server to process the incoming messages, checks if failed to serve
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
