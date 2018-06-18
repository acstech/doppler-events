package couchbase

import (
	"errors"
	"fmt"
	"net/url"
	"sort"

	"github.com/couchbase/gocb"
)

// Client is the client that is to be checked for in the bucket.
// EventID is the a string which is the eventID of this particular event.
// Conn is a connection string which be of the form couchbase://username:password@localhost/bucket_name.
// Bucket is the gocb Bucket of the bucket listed in the connection string.
type Client struct {
	EventID string
	Conn    string
	Bucket  *gocb.Bucket
	doc     *doc
	docID   gocb.Cas
}

// doc is the document that is associated with the client.
type doc struct {
	ID     string   `json:"ID,omitempty"`
	Events []string `json:"Events"`
}

/*func main() {
	c := &Client{EventID: "1", Conn: "couchbase://validator:rotadilav@localhost/doppler"}
	c.SetClientID("1")
	fmt.Println("Connecting to couchbase...")
	c.ConnectToCB()
	fmt.Println("Connected to couchbase.")
	if c.ClientExists() {
		fmt.Println("Client " + c.doc.ID + " exists")
		c.EventEnsure()
	} else {
		fmt.Println("Creating new client")
		c.CreateDocument()
		fmt.Println("Created the document")
	}
	fmt.Println("Closing the client's bucket.")
	c.Bucket.Close()
}*/

// SetClientID ...
func (c *Client) SetClientID(ID string) {

	c.doc.ID = ID
}

// ClientExists determines whether or not a couchbase client exists or not.
// returns true if the document exists and false otherwise.
func (c *Client) ClientExists() bool {
	err := c.collectEvents()
	if err != nil {
		// check to see if the key exists
		if err.Error() == "key not found" {
			fmt.Println("The document was not found")
			return false
		}
		panic(err)
	}
	return true
}

// collectEvents gets the list of eventID's for the client from the couchbase document.
func (c *Client) collectEvents() error {
	var err error
	var docFrag *gocb.DocumentFragment
	docFrag, err = c.Bucket.LookupIn(fmt.Sprintf("doppler:client:%s", c.doc.ID)).Get("Events").Execute()
	if err != nil {
		return err
	}
	// get the document's cas value so that the client data can be updated later if need be
	c.docID = docFrag.Cas()
	// get the Events array into a slice
	docFrag.Content("Events", &c.doc.Events)
	if err != nil {
		return err
	}
	// display the results of the array to slice conversion
	fmt.Printf("%#v\n", c.doc.Events)
	return nil
}

// EventEnsure determines if an event exest for a particular client.
// If it does, it will stop there, but if the event doesn't exist it will add it to the slice in its appropriate place
// and updates the document.
func (c *Client) EventEnsure() {
	if binarySearch(c.doc.Events, c.EventID) == -1 {
		fmt.Println("The eventID does not exist, so add the event to the slice")
		// the event does not exist so add it to and sort the slice
		c.doc.Events = append(c.doc.Events, c.EventID)
		sort.SliceStable(c.doc.Events, func(i, j int) bool { return c.doc.Events[i] < c.doc.Events[j] })
		fmt.Printf("%#v\n", c.doc.Events)
		// update the document
		_, err := c.Bucket.MutateIn(fmt.Sprintf("doppler:client:%s", c.doc.ID), c.docID, 0).Replace("Events", c.doc.Events).Execute()
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("The eventID " + c.EventID + " exists")
	}
}

// CreateDocument creates a document for the client in the couchbase bucket.
func (c *Client) CreateDocument() {
	_, err := c.Bucket.Upsert(fmt.Sprintf("doppler:client:%s", c.doc.ID), c.doc.Events, 0)
	if err != nil {
		panic(err)
	}
}

// ConnectToCB connects to couchbase and sets up the client's bucket.
// returns an error which will be nil if no error has occured.
func (c *Client) ConnectToCB() error {
	// parse the connection string into a url for later use
	u, err := url.Parse(c.Conn)
	if err != nil {
		return err
	}
	// make sure that the url is going to couchbase
	if u.Scheme != "couchbase" {
		return errors.New("Scheme must be couchbase")
	}
	// make sure that a username and password exist, this is required by couchbase 5 and higher
	username, password := "", ""
	if u.User != nil {
		username = u.User.Username()
		password, _ = u.User.Password()
	}
	// make sure that the bucket to connect to is specified
	if u.Path == "" || u.Path == "/" {
		return errors.New("Bucket not specified")
	}
	// get the proper connection format (couchbase//host) and connect to the cluster
	spec := fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	cluster, err := gocb.Connect(spec)
	if err != nil {
		return err
	}
	// authenticate the user and connect to the specified bucket
	cluster.Authenticate(&gocb.PasswordAuthenticator{Username: username, Password: password})
	c.Bucket, err = cluster.OpenBucket(u.Path[1:], "")
	if err != nil {
		return err
	}
	return nil
}

// GetEvents gets the events list and returns them as a slice.
func (c *Client) GetEvents() []string {
	if c.doc.Events == nil {
		c.collectEvents()
	}
	return c.doc.Events
}

// binarySearch impliments a binarySearch on a slice of strings.
// returns -1 if the search string is not found and the index where the search string is in the slice if it is found.
// variation of https://stackoverflow.com/questions/43073681/golang-binary-search
func binarySearch(slice []string, search string) (result int) {
	mid := len(slice) / 2
	switch {
	case len(slice) == 1:
		if slice[0] == search {
			result = 0
		} else {
			result = -1 // not found
		}
	case slice[mid] > search:
		result = binarySearch(slice[:mid], search)
	case slice[mid] < search:
		result = binarySearch(slice[mid+1:], search)
		result += mid + 1
	default: // a[mid] == search
		result = mid // found
	}
	return
}
