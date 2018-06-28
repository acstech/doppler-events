// Couchbase Connector is used for connecting to and working with a couchbase bucket.

package couchbase

import (
	"errors"
	"fmt"
	"net/url"
	"sort"

	"github.com/couchbase/gocb"
)

// Couchbase is a couchbase bucket connection.
// Bucket is the bucket listed in the connection string.
// Doc is the document structure a client document.
type Couchbase struct {
	Bucket     *gocb.Bucket
	bucketName string
}

// Doc is the document that is associated with the client.
// ID is the client's ID.
// Events is the list of client eventID's
type Doc struct {
	ID     string   `json:"ID,omitempty"`
	Events []string `json:"Events"`
}

// ClientExists determines whether or not a couchbase client exists or not.
// clientID is the client's ID.
// returns true if the document exists and false otherwise.
func (c *Couchbase) ClientExists(clientID string) (bool, *Doc, error) {
	document, err := c.collectEvents(clientID)
	if err != nil {
		// check to see if the key exists
		if gocb.IsKeyNotFoundError(err) {
			return false, nil, nil
		}
		return false, nil, err
	}
	return true, document, nil
}

// collectEvents gets the list of eventID's for the client from the couchbase document.
// clientID is the client's ID.
func (c *Couchbase) collectEvents(clientID string) (*Doc, error) {
	var err error
	var docFrag *gocb.DocumentFragment
	document := &Doc{}
	docFrag, err = c.Bucket.LookupIn(fmt.Sprintf("%s:client:%s", c.bucketName, clientID)).Get("Events").Execute()
	if err != nil {
		return nil, err
	}
	// get the Events array into a slice
	docFrag.Content("Events", &document.Events)
	if err != nil {
		return nil, err
	}
	return document, nil
}

// EventEnsure adds the provided event to the client's document if it is not already there.
// clientID is the client's ID.
// eventID is the client's event ID.
func (c *Couchbase) EventEnsure(clientID, eventID string, document *Doc) error {
	sort.Sort(sort.StringSlice(document.Events))
	// determine if the eventID is already in couchbase, if it is not then add it
	if binarySearch(document.Events, eventID) == -1 {
		_, err := c.Bucket.MutateIn(fmt.Sprintf("%s:client:%s", c.bucketName, clientID), 0, 0).ArrayAddUnique("Events", eventID, false).Execute()
		if err != nil {
			if err.Error() != "subdocument mutation 0 failed (given path already exists in the document)" {
				return err
			}
		}
	}
	return nil
}

// ConnectToCB connects to couchbase and sets up the client's bucket.
// conn is a connection string which be of the form couchbase://username:password@localhost/bucket_name.
// returns an error which will be nil if no error has occured.
func (c *Couchbase) ConnectToCB(conn string) error {
	// parse the connection string into a url for later use
	u, err := url.Parse(conn)
	if err != nil {
		return err
	}
	// make sure that the url is going to couchbase
	if u.Scheme != "couchbase" {
		return errors.New("Scheme must be couchbase, verify .env is correct")
	}
	// make sure that a username and password exist, this is required by couchbase 5 and higher
	username, password := "", ""
	if u.User != nil {
		username = u.User.Username()
		password, _ = u.User.Password()
	}
	// make sure that the bucket to connect to is specified
	if u.Path == "" || u.Path == "/" {
		return errors.New("Bucket not specified, verify .env is correct")
	}
	c.bucketName = u.Path[1:]
	// get the proper connection format (couchbase//host) and connect to the cluster
	spec := fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	cluster, err := gocb.Connect(spec)
	if err != nil {
		return err
	}
	// authenticate the user and connect to the specified bucket
	cluster.Authenticate(&gocb.PasswordAuthenticator{Username: username, Password: password})
	c.Bucket, err = cluster.OpenBucket(c.bucketName, "")
	if err != nil {
		return err
	}
	return nil
}

// binarySearch does a binary search and determines if an element exists inside of a slice or not
// modified version of https://stackoverflow.com/questions/43073681/golang-binary-search
func binarySearch(a []string, search string) (result int) {
	mid := len(a) / 2
	switch {
	case len(a) == 0:
		result = -1 // not found
	case a[mid] > search:
		result = binarySearch(a[:mid], search)
	case a[mid] < search:
		result = binarySearch(a[mid+1:], search)
		if result != -1 {
			result += mid + 1
		}
	default: // a[mid] == search
		result = mid // found
	}
	return
}
