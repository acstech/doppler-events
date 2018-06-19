# doppler-events

An API to connect clients to a queue to process for live preview.

In /private/etc/hosts, add "kakfa" next to "localhost" so it can choose between hostnames
To send test event from dummy client to server:

go run cmd/grpcTEST/serviceStart.go

go run cmd/grpcTEST/testsend/testSend.go


## Setup
### Couchbase
- After starting Couchbase through docker, go to localhost:8091/ and create an administrator, choose to configure the memory and choose 256 MB, and then create a bucket.
- Then create a user that only has read and write permissions on that same bucket using the security tab.
- Next copy and rename .envDefualt from data/couchbase/ to the base directory of this repository as .env.
- After that, fill in the appropriate environment varialbe(s).
- Create a document for the client in cmd/grpcTEST/serviceStart.go using the format doppler:client:client_id_here
## Running Development
go run ./cmd/eventAPI/*.go

### Running gRPC Test
go run ./cmd/grpcTEST/serviceStart.go
go run ./cmd/grpcTEST/main.go

