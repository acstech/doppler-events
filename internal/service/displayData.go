package service

import (
	"encoding/json"
	"fmt"
	pb "github.com/acstech/doppler-events/rpc/eventAPI"
	"github.com/couchbase/gocb"
	ptype "github.com/golang/protobuf/ptypes"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"time"
)

// DisplayData is the function that EventAPIClient.go calls in order to send data to the server
// the data is then processed, formatted to JSON, and send to Kafka
func (s *Service) DisplayData(ctx context.Context, in *pb.DisplayRequest) (*pb.DisplayResponse, error) {

	errs := s.verifyConstraints(in)
	if len(errs.errMes) != 0 {
		errorMSG := ""
		for e := range errs.errMes {
			errorMSG += errs.errMes[e] + ", "
		}
		return nil, status.Error(codes.InvalidArgument, "401: Invalid input: "+errorMSG[:len(errorMSG)-2])
	}
	//converting protobuf timestap to to a string in format yyyy-MM-DDTHH:mm:ss.SSSZ
	ts := ptype.TimestampString(in.DateTime)
	//make sure that the timestamp is before now
	now, err := ptype.TimestampProto(time.Now())
	if err != nil {
		return nil, status.Error(codes.Internal, "504: Unable to get the proper time")
	}
	tempTime, err := ptype.Timestamp(in.DateTime)
	if err != nil {
		return nil, status.Error(codes.Internal, "401: Unable to get the proper time")
	}
	if tempTime.After(time.Now()) {
		ts = ptype.TimestampString(now)
	}
	//convert DisplayRequest to map in order to flatten (needed to flatten for influxDB)
	//intialize flatJSONMap as placeholder for marshal
	flatJSONMap := make(map[string]string)
	//check to make sure that the ClientID exists
	cont, document, err := s.CbConn.ClientExists(in.ClientId)
	if err != nil {
		if err == gocb.ErrTimeout {
			return nil, status.Error(codes.Internal, "501: Unable to validate clientID")
		} else if err == gocb.ErrBusy {
			return nil, status.Error(codes.Internal, "502: Unable to validate clientID")
		}
		return nil, status.Error(codes.Internal, "503: Unable to validate clientID")
	}
	if !cont {
		return nil, status.Error(codes.NotFound, "401: The ClientID is not valid")
	}
	//ensure that the eventID exists
	err = s.CbConn.EventEnsure(in.ClientId, in.EventId, document)
	if err != nil {
		//an error ensuring that the event be added to couchbase
		if err == gocb.ErrTimeout {
			return nil, status.Error(codes.Internal, "501: Unable to validate clientID")
		} else if err == gocb.ErrBusy {
			return nil, status.Error(codes.Internal, "502: Unable to validate clientID")
		}
		return nil, status.Error(codes.Internal, "503: Unable to validate clientID")
	}
	//will always have clientID, eventID, dateTime
	flatJSONMap["clientID"] = in.ClientId
	flatJSONMap["eventID"] = in.EventId
	flatJSONMap["dateTime"] = ts
	//loop across dataSet map and add key and value to flatJSON
	for key, value := range in.DataSet {
		if key == "lat" {
			val, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, status.Error(codes.InvalidArgument, "401: Invalid input: lattitude type error")
			}
			if val > 85.0 {
				value = "85.0"
			} else if val < -85.0 {
				value = "-85.0"
			}
		} else if key == "lng" {
			val, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, status.Error(codes.InvalidArgument, "401: Invalid input: lattitude type error")
			}
			if val > 175.0 {
				value = "175.0"
			} else if val < -175.0 {
				value = "-175.0"
			}
		}
		flatJSONMap[key] = value
	}

	//format to JSON
	JSONbytes, err := json.Marshal(flatJSONMap) //Marshal returns the ascii presentation of the data
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "401: Invalid input")
	}
	s.sendToQueue(JSONbytes)
	//return response to client
	return &pb.DisplayResponse{Response: fmt.Sprintf("Success: %s", in.ClientId)}, nil
}
