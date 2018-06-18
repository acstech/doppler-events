package service

import (
	"context"
	"fmt"

	"github.com/dixonwille/golang-template/rpc/serviceName"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	// This is to make sure that your service implements the proper interface.
	_ = serviceName.ServiceNameServer(&Service{})
}

//Service that is very polite
type Service struct{}

//New creates a Service
func New() *Service {
	return &Service{}
}

//SayHello to the person who is requesting conversation
func (*Service) SayHello(_ context.Context, req *serviceName.HelloRequest) (*serviceName.HelloResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "Name is required")
	}
	msg := fmt.Sprintf("Hello, %s", req.Name)
	return &serviceName.HelloResponse{
		Message: msg,
	}, nil
}
