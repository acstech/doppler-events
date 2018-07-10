# Doppler Heatmap
This project provides and API for developers to be able to send events along with a coordinate point and have it plot live on a heatmap hosted in the browser. Written in Go and uses gRPC, Kafka, Couchbase, InfluxDB, Leaflet Maps, and Heatmap.js.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites


* [go](https://golang.org/doc/install) (written for go1.10.3)
* [Docker-compose](https://docs.docker.com/compose/install/#install-compose)


### Installing

Setting up the evironment

In directory of your choice, clone Doppler-Events (Backend API/Producer), Doppler-API (Frontend API/Consumer), and Doppler-Frontend (Frontend)


```
git clone https://github.com/acstech/doppler-events.git
git clone https://github.com/acstech/doppler-api.git
git clone https://github.com/acstech/doppler-frontend.git
```

**Environment Variables**

Rename the doppler-base.env.default file to doppler-base.env and fill out the environment variables by following the comments.

_Note: anything with an underscore in the name should be replaced with the appropriate value (ip and port need to be changed as well)._

**Docker Setup**

In doppler-events directory:

Run `docker-compose up -d`

Run `docker-compose ps` to make sure all services are up except doppler-api and doppler-events because couchbase has not yet been setup.

Example output:

```
Name                                  Command                      State     Ports                         
----------------------------------------------------------------------------------------------------------
doppler-events_doppler-api_1       ./entrypoint.sh                  Up      0.0.0.0:8000->8000/tcp
```

**Kafka Setup**

Map "Kafka" to localhost in your host file

    Mac: Open /private/etc/hosts in a text editor and add 'kafka' (no quotes) after
    127.0.0.1	localhost

	Linux: TODO

    Windows: TODO

**Couchbase Setup**

After docker is finished starting up Couchbase, visit the host:port, specified in the docker-compose.yml, file in your browser to set up Couchbase (the default host and port are localhost and 8091)

_Note: Couchbase can take up to a minute after docker is running to be served to the specified host._

Create an account with a username and password of your choice.

Go to Buckets and Add Bucket (Top right)

Create a bucket (save the bucket name you give it somewhere) and configure the memory to your liking. 256 MB is recommended.

Create a user that only had read and write permissions on that same bucket using the security tab. (Located under Data Roles->Data Writer and Data Roles-> Data Reader)

Create a document and name it with the following format without brackets:

    [YourBucketName]:client:[YourClientName]

In the JSON for the document, add the client name and an empty Events array and save. Example:

    {
    	"ID":"YourClientName",
        "Events":[]
    }

Make sure to run `docker-compose up -d` again becuase doppler-api and doppler-events crashed due to not being able to connect to couchbase.

## Running the tests

At this point, if you have had no error messages come, everything is standing up. Visit the location you served the front-end to (127.0.0.1:9080 by default) in your browser and enter in your clientID. Settings can be found by clicking the hamburger icon on the top left.

Now lets try sending it some test data. Edit the doppler-events/cmd/testsend/testSend.go on line 29 and 30 add your client and eventIDs. Then run the testSend.go file: go run testSend.go. You can now go to your map and view the test data send.

## Local Development and Testing

Whenever a file changes inside of any of the three repos (doppler-api, doppler-events, doppler-frontend), the appropriate docker build command will need to be run.

In doppler-events run: docker build . -t acstintern/doppler-events:latest

In doppler-api run: docker build . -t acstintern/doppler-api:latest

In doppler-frontend run: docker build . -t acstintern/doppler-frontend:latest

And make sure that the docker-compose.yml file has either ':' or ':latest' after the doppler services, which correspond to the different repositories.

Then run `docker-compose up -d` and proceed to checkout the new changes.

To send data all the way through, at this point, run `go run cmd/testsend/testSend.go`.
To run testSend.go with flags, use `-s` to simulate data for America only, `-p` to hit geographic locations multiple times, `-l` for infinite random latitude and longitude points. This test file has database cleanup for influx implemented after the test is finished.

To load test the API, run `go run cmd/loadtest/loadTest.go`. This will not clean up your database.

## Deployment

Coming soon

## Built With

* [go](https://golang.org/) - Backend
* [kafka](http://kafka.apache.org/) - Messaging Queue
* [influx-sink](https://lenses.stream/connectors/sink/influx.html) - Used to sink data from kafka into Influx
* [influxDB](https://www.influxdata.com/) - Used to store time based data
* [Couchbase](https://www.couchbase.com/) - Used to store clients and events
* [jQuery](https://jquery.com/) - Front End
* [Bootstrap](https://getbootstrap.com/) - Frontend css
* [Leaflet](https://leafletjs.com/index.html) - Frontend heatmap

## Contributing

Coming soon

## Versioning

Version 1.0.0

## Authors

* ***Matt Harrington***
* ***Peter Kaufmann***
* ***Pranav Minasandram*** - gRPC and Frontend
* ***Matt Smith***
* ***Leander Stevano***
* ***Ben Wornom***

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.

## License

## Acknowledgments

* [Heatmap.js](https://www.patrick-wied.at/static/heatmapjs/) by Patrick Wied
