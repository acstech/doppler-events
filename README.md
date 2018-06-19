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
- Next add `export COUCHBASE_HOST="couchbase_host"
export COUCHBASE_BUCKET="bucket_name"
export COUCHBASE_USERNAME="couchbase_bucket_user"
export COUCHBASE_PASSWORD="couchbase_bucket_user_password"` to your bash startup script.
- Create a document for the client in cmd/grpcTEST/serviceStart.go using the format doppler:client:client_id_here
## Running Development
go run ./cmd/eventAPI/*.go

### Running gRPC Test
go run ./cmd/grpcTEST/serviceStart.go
go run ./cmd/grpcTEST/main.go

