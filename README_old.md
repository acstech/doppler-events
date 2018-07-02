# Doppler-Events
An API to connect clients to a queue to process for live preview.
## Setup
### General
- Clone the repository
- Install [go](https://golang.org/dl/) and [docker](https://docs.docker.com/install/)
- Run `docker-compose up -d`
- Map "kafka" to localhost in your host file
### Couchbase
- Go to localhost:8091 (if Couchbase does not startup after a minute run `docker-compose down` and try this step again)
- Create an administrator
- Choose to configure the memory and choose 256 MB
- Create a bucket
- Create a user that only has read and write permissions on that same bucket using the security tab (these permissions are located under Data Roles->Data Writer and Data Roles->Data Reader)
- Next copy and rename .env.defualt from data/couchbase/ to the base directory of this repository as .env
- After that, fill in the appropriate environment varialbe
## Testing **** Not completed yet ****
- Create a documents for each client that you are going to use in the aformentioned bucket following the format bucket_name:client:clientID
Note:  client is going to be the same no matter what the client name is
- Run `go run cmd/grpcTEST/serviceStart.go`
- Run `go run cmd/testsend/testSend.go` (The clientID for this test are client0 and nav1)
- Check influx by running `docker-compose exec influx_db influx`
- Next type `USE doppler` and then `SELECT * FROM dopplerDataHistory`
- You should see data in the database for however many times the testSend.go file ran