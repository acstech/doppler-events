package main

import (
	"log"
//	"os"
	"time"
	"strconv"
	"flag"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "github.com/acstech/doppler-events/eventAPI"
)

const (
	address = "localhost:8080"
	//address     = "10.22.97.107:8080"
	defaultName = "default"
)

func main() {
	clientID := flag.String("cid", "N/A", "Client ID")
	eventID := flag.String("eid", "N/A", "Event ID")
	lon := flag.Float64("lon", 0, "Longitude")
	lat := flag.Float64("lat", 0, "Latitude")

	flag.Parse()

	coordinates := make(map[string]string)
	coordinates["lon"] = strconv.FormatFloat(*lon, 'E', -1, 64)
	coordinates["lat"] = strconv.FormatFloat(*lat, 'E', -1, 64)

  conn, err := grpc.Dial(address, grpc.WithInsecure())
  if err != nil{
    log.Fatalf("Did not connect: %v", err)
  }

  defer conn.Close()
  c:= pb.NewEventSenderClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SendEvent(ctx, &pb.EventObj{ClientID: *clientID, EventID: *eventID, TimeSinceEpoch: uint64(time.Now().Unix()), KeyValues: coordinates})
	if err != nil {
		log.Fatalf("could not do anything: %v", err)
	}
	log.Printf(r.Response)
}
