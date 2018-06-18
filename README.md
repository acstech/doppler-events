# doppler-events

An API to connect clients to a queue to process for live preview.

In /private/etc/hosts, add "kakfa" next to "localhost" so it can choose between hostnames
To send test event from dummy client to server:

go run cmd/grpcTEST/serviceStart.go

go run cmd/grpcTEST/testsend/testSend.go


## Setup
### Couchbase
After starting Couchbase through docker, go to localhost:8091/ and create an administrator, choose to configure the memory and choose 1000 MB, and then create a user

## Running Development
go run ./cmd/eventAPI/*.go

### Running gRPC Test
go run ./cmd/grpcTEST/serviceStart.go
go run ./cmd/grpcTEST/main.go

### Running testprod
go run ./cmd/testprod/*.go
