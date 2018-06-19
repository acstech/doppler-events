package main

import (
	"os"

	"github.com/acstech/doppler-events/internal/service"
)

func main() {
	//get config variables
	cbHost := os.Getenv("COUCHBASE_HOST")
	cbBucket := os.Getenv("COUCHBASE_BUCKET")
	cbUser := os.Getenv("COUCHBASE_USERNAME")
	cbPass := os.Getenv("COUCHBASE_PASSWORD")
	//pass config variables so that they can be used later
	service.Init(cbHost, cbBucket, cbUser, cbPass)

}
