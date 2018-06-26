// Event API Server
//

package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	cb "github.com/acstech/doppler-events/internal/couchbase"
	pb "github.com/acstech/doppler-events/rpc/eventAPI"
	ptype "github.com/golang/protobuf/ptypes"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

//intialize addressess of where server listens with format IP:Port
const (
	address = ":8080"
)

type ErrorRes struct {
	err    error
	errMes []string
}

//struct to hold parameters for server
type server struct {
	theProd sarama.AsyncProducer
	cbConn  *cb.Couchbase
}

// New returns an error that formats as the given text.
func New(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

//verifyConstraints checks the attributes for in incoming request, verfies valid data
func verifyConstraints(req *pb.DisplayRequest) ErrorRes {
	var errRes ErrorRes
	//check length of EventId
	if len(req.EventId) > 35 {
		errRes.errMes = append(errRes.errMes, "EventId must be less than 35 characters")
	}
	// check length of ClientId
	if len(req.ClientId) == 0 {
		errRes.errMes = append(errRes.errMes, "ClientId must be included")
	}
	//check for longitude and latitude keys in DataSet
	for key, value := range req.DataSet {
		if key == "lat" {
			//verify "lat" is in proper format
			if reflect.TypeOf(value) == reflect.TypeOf("string") {
				floater, err := strconv.ParseFloat(value, 64)
				if err != nil {
					errRes.errMes = append(errRes.errMes, "invalid latitude value. check type")
				}
				//check valid ranges
				if floater < -85 || floater > 85 {
					errRes.errMes = append(errRes.errMes, "invalid latitude value. check range")
				}
			} else {
				errRes.errMes = append(errRes.errMes, "latitude needs to be a string")
			}
			//verify "lng" is in proper format
		} else if key == "lng" {
			if reflect.TypeOf(value) == reflect.TypeOf("string") {
				floater, err := strconv.ParseFloat(value, 64)
				if err != nil {
					errRes.errMes = append(errRes.errMes, "invalid longitude value. check type")
				}
				if floater < -175 || floater > 175 {
					errRes.errMes = append(errRes.errMes, "invalid longitude value. check range")
				}
			} else {
				errRes.errMes = append(errRes.errMes, "latitude needs to be a string")
			}
			// might want to allow for event to be passed without longitude or latitude
		} else {
			errRes.errMes = append(errRes.errMes, "could not find latitude or longitude")
		}
	}
	return errRes
}

//newProducer configures an asynchronous kafka producer client, returns it
func newProducer() (sarama.AsyncProducer, error) {
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	// Setup configuration
	config := sarama.NewConfig()
	config.ClientID = "1"
	//configuration for batches
	config.Producer.Flush.MaxMessages = 30 // will flush if 30 messages arrived
	config.Producer.Flush.Frequency = 50 * time.Millisecond
	config.Producer.Flush.Messages = 1 // can flush with 1 message
	//The level of acknowledgement reliability needed from the broker.
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// The level of acknowledgement reliability needed from the broker.
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
	// var enqueued, errors int

	msg := &sarama.ProducerMessage{
		Topic: "influx-topic",
		Value: sarama.ByteEncoder(JSONob),
		//Partition: 1, //partNum,
	}
	go func() {
		for err := range prod.theProd.Errors() {
			fmt.Println("Failed to produce message:", err)
		}
	}()
	prod.theProd.Input() <- msg

	// log.Printf("Enqueued: %d; errors: %d\n", enqueued, errors)
}

// DisplayData is the function that EventAPIClient.go calls in order to send data to the server
// the data is then processed, formatted to JSON, and send to Kafka
func (s *server) DisplayData(ctx context.Context, in *pb.DisplayRequest) (*pb.DisplayResponse, error) {

	errs := verifyConstraints(in)
	if len(errs.errMes) != 0 {
		for e := range errs.errMes {
			fmt.Println("unable to send message:", errs.errMes[e])
		}
		return nil, errors.New("invalid input")
	}
	//converting protobuf timestap to to a string in format yyyy-MM-DDTHH:mm:ss.SSSZ
	ts := ptype.TimestampString(in.DateTime)

	//convert DisplayRequest to map in order to flatten (needed to flatten for influxDB)
	//intialize flatJSONMap as placeholder for marshal
	flatJSONMap := make(map[string]string)
	//check to make sure that the ClientID exists
	cont, err := s.cbConn.ClientExists(in.ClientId)
	if err != nil {
		return nil, err
	}
	if !cont {
		return nil, errors.New("the ClientID is not valid")
	}
	//ensure that the eventID exists
	s.cbConn.EventEnsure(in.ClientId, in.EventId)
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
		return nil, err
	}
	s.sendToQueue(JSONbytes)
	//return response to client
	return &pb.DisplayResponse{Response: fmt.Sprintf("Success: %s", in.ClientId)}, nil
}

// Init sets up the backend server so that clients can send data to kafka
// cbCon is the connection string for couchbase
// returns an error if any occur while creating a kafka producer, a couchbase connection, sending data,
// or closing the kafka producer
func Init(cbCon string) error {
	//initialize listener on server address
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	//initialize server
	s := grpc.NewServer()
	prod, err := newProducer()
	if err != nil {
		fmt.Println("failed to create Kafka producer connection. Ensure docker is running and the kafka topic is connected.")
		return err
	}

	defer func() error {
		if err := prod.Close(); err != nil {
			// Should not reach here
			return err
		}
		return nil
	}()

	serve2 := server{
		theProd: prod,
		cbConn:  &cb.Couchbase{Doc: &cb.Doc{}},
	}
	err = serve2.cbConn.ConnectToCB(cbCon)
	if err != nil {
		return fmt.Errorf("CB connection error: %v", err)
	}
	//register server to grpc
	pb.RegisterEventAPIServer(s, &serve2)
	// Register reflection service on gRPC server for back and forth communication. Kept for future use if necessary
	// reflection.Register(s)

	//tells the server to process the incoming messages, checks if failed to serve
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
