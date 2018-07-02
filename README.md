# Doppler Heatmap
This project provides and API for developers to be able to send events along with a coordinate point and have it plot live on a heatmap hosted in the browser. Written in Go and uses gRPC, Kafka, Couchbase, InfluxDB, Leaflet Maps, and Heatmap.js.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites


* [go](https://golang.org/doc/install) (written for go1.10.3)

* [Docker-compose](https://docs.docker.com/compose/install/#install-compose)

* live website server of choice ([npm static-server](https://www.npmjs.com/package/static-server), atom-live-server, etc.)

***Optional***

* [npm](https://www.npmjs.com/get-npm) (Only if you are using static-server)


### Installing

Setting up the evironment

In directory of your choice, clone Docker-Events (Backend API/Producer), Docker-API (Frontend API/Consumer), and Docker-Frontend (Frontend)


```
git clone https://github.com/acstech/doppler-events.git
git clone https://github.com/acstech/doppler-api.git
git clone https://github.com/acstech/doppler-frontend.git
```

In doppler-events directory:
```
docker-compose up -d
```
**Kafka Setup**

Map "Kafka" to localhost in your host file
	
    Mac: Open /private/etc/hosts in a text editor and add 'kafka' (no quotes) after
    127.0.0.1	localhost

	Linux: TODO
    
    Windows: TODO

**Couchbase Setup**

After docker is finished installing, visit localhost:8091 in your browser to set up Couchbase

_Note: Couchbase can take up to a minute after docker is running to be served to localhost._

Create an account with a username and password of your choice.

Go to Buckets and Add Bucket (Top right)

Create a bucket (save the bucket name you give it somewhere) and configure the memory to your liking. 256 MB is recommended.

Create a client that only had read and write permissions on that same bucket using the security tab. (Located under Data Roles->Data Writer and Data Roles-> Data Reader)

Create a document and name it with the following format without brackets:
	
    [YourBucketName]:client:[YourClientName]

In the JSON for the document, add the client name and an empty Events array and save. Example:
	
    {
    	"ID":"YourClientName",
        "Events":[]
    }
    
In the doppler-events/data/couchbase/ directory, copy the .env.default file and put it in the base directory (doppler-events)

Remove “.default” from the filename.

Edit the file and replace fields with appropriate information

Example:
```
### Couchbase
export COUCHBASE_CONN="couchbase://ExampleUserName:topS3cretPassword123@localhost/YourBucket”
```

Copy the same .env file and place it in the root of the doppler-api directory

**Frontend Setup**

Serve your website using live-server of your choice.

Example using static-server on npm (Must have [npm](https://www.npmjs.com/get-npm) installed):

    To install static-server: npm -g install static-server
	In the doppler-frontend directory, run the command: static-server

**Start backend API**

From the doppler-events directory:
		
   	go run cmd/grpcTEST/serviceStart.go

**Start frontend api**

From the doppler-api directory:

	go run cmd/doppler-api/main.go

***Note: It is important that the last two commands are run from their respective base directories, as is.*** 




## Running the tests

At this point, if you have had no error messages come up, everything is standing up. Visit the location you served the front-end to (127.0.0.1:9080 if you used static-server) in your browser and enter in your clientID. Settings can be found by clicking the hamburger icon on the top left.

Now lets try sending it some test data. Edit the doppler-events/cmd/testsend/testSend.go on line 29 and 30 add your client and eventIDs. Then run the testSend.go file. go run testSend.go. You can now go to your map and view the test data send. 

## Deployment

TODO

## Built With

* [go](https://golang.org/) - Backend
* [kafka](http://kafka.apache.org/) - Messaging Queue
* [influxDB](https://www.influxdata.com/) - Used to store time based data
* [Couchbase](https://www.couchbase.com/) - Used to store clients and events
* [jQuery](https://jquery.com/) - Front End

## Contributing


TODO

Please read [CONTRIBUTING.md](https://gist.github.com/PurpleBooth/b24679402957c63ec426) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

TODO

## Authors

* ***Matt Harrington***
* ***Peter Kaufmann***
* ***Pranav Minasandram*** - gRPC and Frontend
* ***Matt Smith***
* ***Leander Stevano***
* ***Ben Wornom***

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.

## License

TODO

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* [Heatmap.js](https://www.patrick-wied.at/static/heatmapjs/) by Patrick Wied
