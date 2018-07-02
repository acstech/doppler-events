package main

import (
	"os"

	"github.com/acstech/doppler-events/internal/service"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	fmt.Println("Error loading .env file")
	// 	panic(err)
	// }
	//get config variables
	cbConn := os.Getenv("COUCHBASE_CONN")
	//pass config variables so that they can be used later
	err := service.Init(cbConn)
	// if an error occurred in itialization, shutdown the server
	if err != nil {
		panic(err)
	}
}
