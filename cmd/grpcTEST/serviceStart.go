package main

import (
	"github.com/acstech/doppler-events/internal/service"
	//c meaning client call
)
const{
	var COUCHBASE_HOST string = "localhost"
	var COUCHBASE_BUCKET string = "doppler"
	var COUCHBASE_USERNAME string = "validator"
	var COUCHBASE_PASSWORD string = "rotadilav"
}
func main() {
	
	service.Init(COUCHBASE_HOST, COUCHBASE_BUCKET, COUCHBASE_USERNAME, COUCHBASE_PASSWORD)

}
