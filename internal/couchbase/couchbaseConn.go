// Couchbase Connector is used for connecting to and working with a couchbase bucket.

package couchbase

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/couchbase/gocb"
)

// Couchbase is a couchbase bucket connection.
// Bucket is the bucket listed in the connection string.
// Doc is the document structure a client document.
type Couchbase struct {
	Bucket     *gocb.Bucket
	Doc        *Doc
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
func (c *Couchbase) ClientExists(clientID string) (bool, error) {
	err := c.collectEvents(clientID)
	if err != nil {
		// check to see if the key exists
		if gocb.IsKeyNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// collectEvents gets the list of eventID's for the client from the couchbase document.
// clientID is the client's ID.
func (c *Couchbase) collectEvents(clientID string) error {
	var err error
	var docFrag *gocb.DocumentFragment
	docFrag, err = c.Bucket.LookupIn(fmt.Sprintf("%s:client:%s", c.bucketName, clientID)).Get("Events").Execute()
	if err != nil {
		return err
	}
	// get the Events array into a slice
	docFrag.Content("Events", &c.Doc.Events)
	if err != nil {
		return err
	}
	return nil
}

// EventEnsure adds the provided event to the client's document if it is not already there.
// clientID is the client's ID.
// eventID is the client's event ID.
func (c *Couchbase) EventEnsure(clientID, eventID string) error {
	_, err := c.Bucket.MutateIn(fmt.Sprintf("%s:client:%s", c.bucketName, clientID), 0, 0).ArrayAddUnique("Events", eventID, false).Execute()
	if err != nil {
		if err.Error() != "subdocument mutation 0 failed (given path already exists in the document)" {
			return err
		}
	}
	return nil
}

// CreateDocument creates a document for the client in the couchbase bucket.
// clientID is the client's ID.
// eventID is the client's event ID.
// Note: this resets the Doc data of the connector.
func (c *Couchbase) CreateDocument(clientID, eventID string) error {
	c.Doc.Events = append(c.Doc.Events, eventID)
	_, err := c.Bucket.Upsert(fmt.Sprintf("%s:client:%s", c.bucketName, clientID), c.Doc.Events, 0)
	c.Doc = &Doc{}
	if err != nil {
		return err
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
