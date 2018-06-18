package main

import (
	"log"
	"net"

	"github.com/acstech/doppler-events/internal/service"
	"github.com/dixonwille/golang-template/rpc/serviceName"
	"google.golang.org/grpc"
)

func main() {
	parseConfig()
	lis, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	serviceName.RegisterServiceNameServer(grpcServer, service.New())
	log.Println("Serving on " + address)
	log.Fatalln(grpcServer.Serve(lis))
}
