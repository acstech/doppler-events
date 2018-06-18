//go:generate protoc -I . --go_out=plugins=grpc:. eventAPI.proto

//Package eventAPI is an example of a generated go server and client for gRPC
//You can put documentation on how to use it here as it will show up in godocs.
package eventAPI
