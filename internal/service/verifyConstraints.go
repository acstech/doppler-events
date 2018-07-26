package service

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	pb "github.com/acstech/doppler-events/rpc/eventAPI"
	"github.com/microcosm-cc/bluemonday"
)

//verifyConstraints checks the attributes for in incoming request, verfies valid data
func (*Service) verifyConstraints(req *pb.DisplayRequest) ErrorRes {
	var errRes ErrorRes
	//check length of EventId
	// Sanitizing eventID
	bm := bluemonday.UGCPolicy()
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		fmt.Errorf("Regex Compile Error: %v", err)
	}
	req.EventId = reg.ReplaceAllString(bm.Sanitize(req.EventId), "")
	if len(req.EventId) > 35 {
		errRes.errMes = append(errRes.errMes, "EventId must be less than 35 characters")
	}
	//check length of ClientId
	if len(req.ClientId) == 0 {
		errRes.errMes = append(errRes.errMes, "ClientId must be included")
	}
	// check to make sure that lat and lng exist based on size of the dataSet
	if len(req.DataSet) < 2 {
		errRes.errMes = append(errRes.errMes, "could not find latitude or longitude")
	}
	//check for longitude and latitude keys in DataSet
	for key, value := range req.DataSet {
		if key == "lat" {
			//verify "lat" is in proper format
			if reflect.TypeOf(value) == reflect.TypeOf("string") {
				floater, err := strconv.ParseFloat(value, 64)
				if err != nil {
					errRes.errMes = append(errRes.errMes, "invalid latitude type")
				}
				//check valid ranges
				if floater < -90 || floater > 90 {
					errRes.errMes = append(errRes.errMes, "invalid latitude value")
				}
			} else {
				errRes.errMes = append(errRes.errMes, "latitude needs to be a string")
			}
			//verify "lng" is in proper format
		} else if key == "lng" {
			if reflect.TypeOf(value) == reflect.TypeOf("string") {
				floater, err := strconv.ParseFloat(value, 64)
				if err != nil {
					errRes.errMes = append(errRes.errMes, "invalid longitude type")
				}
				if floater < -180 || floater > 180 {
					errRes.errMes = append(errRes.errMes, "invalid longitude value")
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
