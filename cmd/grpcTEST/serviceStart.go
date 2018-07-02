package main

import (
<<<<<<< HEAD
	"errors"
	"net/url"
=======
>>>>>>> Updated docker images
	"os"

	"github.com/acstech/doppler-events/internal/service"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
<<<<<<< HEAD
	//get config variables
	cbConn := os.Getenv("COUCHBASE_CONN")
	address := os.Getenv("API_ADDRESS")
	kafkaConn, kafkaTopic, err := kafkaParse(os.Getenv("KAFKA_CONN"))
	if err != nil {
		panic(err)
	}
	//pass config variables so that they can be used later
	err = service.Init(cbConn, kafkaConn, kafkaTopic, address)
=======
	// err := godotenv.Load()
	// if err != nil {
	// 	fmt.Println("Error loading .env file")
	// 	panic(err)
	// }
	//get config variables
	cbConn := os.Getenv("COUCHBASE_CONN")
	//pass config variables so that they can be used later
	err := service.Init(cbConn)
>>>>>>> Updated docker images
	// if an error occurred in itialization, shutdown the server
	if err != nil {
		panic(err)
	}
}

func kafkaParse(conn string) (string, string, error) {
	u, err := url.Parse(conn)
	if err != nil {
		return "", "", err
	}
	if u.Host == "" {
		return "", "", errors.New("Kafka address is not specified, verify that your environment varaibles are correct")
	}
	address := u.Host
	// make sure that the topic is specified
	if u.Path == "" || u.Path == "/" {
		return "", "", errors.New("Kafka topic is not specified, verify that your environment varaibles are correct")
	}
	topic := u.Path[1:]
	return address, topic, nil
}
