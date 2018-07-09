package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	pb "github.com/acstech/doppler-events/rpc/eventAPI" //c meaning client call
	"github.com/golang/protobuf/ptypes"
	client "github.com/influxdata/influxdb/client/v2"
	"google.golang.org/grpc"
)

var (
	clientIDs []string
	eventIDs  []string
	c         pb.EventAPIClient
	stop      bool
)

//Northeast generates points for the northeast USA.
func Northeast() (float64, float64) {
	lat := 35 + rand.Float64()*(44-35)
	lng := 74 + rand.Float64()*(83-74)
	lng = lng - lng - lng
	return lat, lng
}

//Southeast usa
func Southeast() (float64, float64) {
	lat := 31 + rand.Float64()*(35-31)
	lng := 79 + rand.Float64()*(87-79)
	lng = lng - lng - lng
	return lat, lng
}

//Midwest USA points
func Midwest() (float64, float64) {
	lat := 29 + rand.Float64()*(48-29)
	lng := 90 + rand.Float64()*(108-90)
	lng = lng - lng - lng
	return lat, lng
}

//Western USA points
func West() (float64, float64) {
	lat := 33 + rand.Float64()*(48-33)
	lng := 108 + rand.Float64()*(121-108)
	lng = lng - lng - lng
	return lat, lng
}

//Random point generator (within USA)
func Random() (float64, float64) {
	lat := 26 + rand.Float64()*(48-26)
	lng := 80 + rand.Float64()*(118-80)
	lng = lng - lng - lng
	return lat, lng
}

//Simulate creates data points that simulates activity sweeping across America from East to West coast
func Simulate() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var lng, lat float64
	for count := 0; count < 1500; count++ {
		time.Sleep(10 * time.Millisecond)
		clientID := clientIDs[r.Int31n(int32(len(clientIDs)))] //pick random client from clientIDs slice
		eventID := eventIDs[r.Int31n(int32(len(eventIDs)))]    //pick random event from eventIDs slice
		if count < 200 {
			lat, lng = Northeast()
		} else if count < 400 {
			if count%5 == 0 {
				lat, lng = Northeast()
			} else {
				lat, lng = Southeast()
			}
		} else if count < 600 {
			if count%5 == 0 {
				lat, lng = Southeast()
			} else if count%10 == 0 {
				lat, lng = Northeast()
			} else {
				lat, lng = Midwest()
			}
		} else if count < 1000 {
			if count%5 == 0 {
				lat, lng = Midwest()
			} else if count%10 == 0 {
				lat, lng = Southeast()
			} else if count%15 == 0 {
				lat, lng = Northeast()
			} else {
				lat, lng = West()
			}
		} else {
			lat, lng = Random()
		}
		resp, err := sendRequest(c, clientID, eventID, lat, lng) //call function that prepares data to send to server
		//lo any error
		if err != nil {
			log.Println(err)
			continue
		}
		if resp != nil {
			//print server response
			log.Println(resp.Response)
		}
	}
}

//Repeat creates a set number of point locations then iterates through them, rehitting an area multiple times
func Repeat() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	clientID := clientIDs[r.Int31n(int32(len(clientIDs)))] //pick random client from clientIDs slice
	eventID := eventIDs[r.Int31n(int32(len(eventIDs)))]    //pick random event from eventIDs slice
	var locations []map[string]string
	for x := 0; x < 200; x++ {
		lat := (r.Float64() - .5) * 180 //get random lat
		lng := (r.Float64() - .5) * 360 //get random lng
		dataSet := make(map[string]string, 2)
		dataSet["lat"] = strconv.FormatFloat(lat, 'g', -1, 64)
		dataSet["lng"] = strconv.FormatFloat(lng, 'g', -1, 64)
		locations = append(locations, dataSet)
	}
	for y := 0; y < 1500; y++ {
		a := rand.Intn(200)
		//get current time
		dateTime := ptypes.TimestampNow()
		//send data to server, returns response and error
		if !stop {
			_, err := c.DisplayData(context.Background(), &pb.DisplayRequest{
				ClientId: clientID,
				EventId:  eventID,
				DateTime: dateTime,
				DataSet:  locations[a],
			})
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

//LoadTest sends infinite random points to the API
func LoadTest() {
	//get true random
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		//time.Sleep(3 * time.Millisecond)
		clientID := clientIDs[r.Int31n(int32(len(clientIDs)))]   //pick random client from clientIDs slice
		eventID := eventIDs[r.Int31n(int32(len(eventIDs)))]      //pick random event from eventIDs slice
		lat := (r.Float64() - .5) * 180                          //get random lat
		lng := (r.Float64() - .5) * 360                          //get random lng
		resp, err := sendRequest(c, clientID, eventID, lat, lng) //call function that prepares data to send to server
		//lo any error
		if err != nil {
			log.Println(err)
			continue
		}
		//print server response
		if resp != nil {
			log.Println(resp.Response)
		}
	}
}

//CleanupInflux will go through influx and delete test data points
func CleanupInflux(theTime int64) {
	//	influxCon := os.Getenv("INFLUX_CONN")
	// creates influx client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: "username",
		Password: "password",
	})
	if err != nil {
		panic(fmt.Errorf("error connecting to influx: %v", err))
	}
	defer c.Close()

	time.Sleep(2 * time.Second)
	curTime := time.Now().UnixNano()

	now := strconv.FormatInt(curTime, 10)
	inTime := strconv.FormatInt(theTime, 10)

	fmt.Printf("delete from dopplerDataHistory where time > %s and time < %s", inTime, now)
	q := client.NewQuery(fmt.Sprintf("delete from dopplerDataHistory where time > %s and time < %s", inTime, now), "doppler", "ns")

	if _, err := c.Query(q); err != nil {
		fmt.Println(err)
	}
}
func main() {
	args := os.Args[1:]
	cleanup := true
	//data point variables
	clientIDs = []string{"client0", "client1"} //In order for test to work, couchbase must contain all 3 clients
	eventIDs = []string{"physical check in", "mobile login", "rest"}
	var err error
	//connect to server
	c, err = dial("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	//var wg sync.WaitGroup
	//wg.Add(1)
	//listens for interrupt, gracefully cleans up
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	startTime := time.Now().UnixNano()
	time.Sleep(500 * time.Millisecond)

	go func() {
		for {
			<-sigs
			stop = true
			CleanupInflux(startTime)
			os.Exit(0)
		}
	}()

	if len(args) == 0 {
		fmt.Println("usage: testsend.go -l [load test] -s [simulation test] -p [repeat point test] -d [no database cleanup]")
	} else {
		for a := 0; a < len(args); a++ {
			if args[a] == "-l" {
				fmt.Println("starting load test...")
				LoadTest()
			}
			if args[a] == "-s" {
				fmt.Println("starting simulation test...")
				Simulate()
			}
			if args[a] == "-p" {
				fmt.Println("starting repeat point test...")
				Repeat()
			}
			if args[a] == "-d" {
				cleanup = false
			}
		}
		if cleanup {
			CleanupInflux(startTime)
		}
	}
	// go func() {
	// 	for {
	// 		select {
	// 		case <-sigs:
	// 			CleanupInflux(startTime)
	// 			os.Exit(0)
	// 		default:
	// 			if len(args) == 0 {
	// 				fmt.Println("usage: testsend.go -l [load test] -s [simulation test] -p [repeat point test] -d [no database cleanup]")
	// 			} else {
	// 				for a := 0; a < len(args); a++ {
	// 					if args[a] == "-l" {
	// 						fmt.Println("starting load test...")
	// 						LoadTest()
	// 					}
	// 					if args[a] == "-s" {
	// 						fmt.Println("starting simulation test...")
	// 						Simulate()
	// 					}
	// 					if args[a] == "-p" {
	// 						fmt.Println("starting repeat point test...")
	// 						Repeat()
	// 					}
	// 					if args[a] == "-d" {
	// 						cleanup = false
	// 					}
	// 				}
	// 			}
	// 			if cleanup {
	// 				CleanupInflux(startTime)
	// 			}
	// 		}
	// 	}
	// }()
	// wg.Wait()
}

//takes client data, sends it to over connection
func sendRequest(c pb.EventAPIClient, clientID, eventID string, lat, lng float64) (*pb.DisplayResponse, error) {
	//create map of data
	dataSet := make(map[string]string, 2)
	dataSet["lat"] = strconv.FormatFloat(lat, 'g', -1, 64)
	dataSet["lng"] = strconv.FormatFloat(lng, 'g', -1, 64)
	//get current time
	dateTime := ptypes.TimestampNow()
	//send data to server, returns response and error

	if !stop {
		resp, err := c.DisplayData(context.Background(), &pb.DisplayRequest{
			ClientId: clientID,
			EventId:  eventID,
			DateTime: dateTime,
			DataSet:  dataSet,
		})
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
	return nil, nil
}

//get grpc connection client
func dial(addr string) (pb.EventAPIClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure()) //WithInsecure meaning no authentication required
	if err != nil {
		return nil, fmt.Errorf("Did not connect: %v", err)
	}
	client := pb.NewEventAPIClient(conn)
	return client, nil
}
