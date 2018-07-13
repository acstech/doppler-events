package main

import (
	"fmt"
	"log"

	client "github.com/influxdata/influxdb/client/v2"
	//"github.com/influxdata/influxdb/client/v2"
)

func filterByEvent(eventid string) error {
	return nil
}

type Response struct {
	Tags   []string
	Values [][]interface{}
}

const (
	DB      = "doppler"
	Measure = "dopplerDataHistory"
)

//GetPoints recieves a query string and returns the results from the influx client and returns the results
func GetPoints(clnt client.Client, query string) (r1 Response, err error) {

	q := client.NewQuery(query, DB, "s")
	response, err := clnt.Query(q)
	if err == nil && response.Error() == nil {
		//	fmt.Println(response.Results)

		for r := range response.Results {
			res := response.Results[r].Series[r].Values
			b := response.Results[r].Series[r].Columns
			r1.Tags = b
			r1.Values = res
			return r1, nil

		}
	}
	return r1, err

}

func main() {

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: "username",
		Password: "password",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	res, err := GetPoints(c, "select lat, lng from dopplerDataHistory")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res.Tags)
		for a := range res.Values {
			res2 := res.Values[a] //    [][]interface {}
			fmt.Println(res2)
		}
	}
}
