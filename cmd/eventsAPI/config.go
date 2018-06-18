package main

import (
	"fmt"

	"github.com/ianschenck/envflag"
	_ "github.com/joho/godotenv/autoload"
)

var (
	port    int
	address string
)

func init() {
	envflag.IntVar(&port, "PORT", 6789, "")
}

func parseConfig() {
	envflag.Parse()
	address = fmt.Sprintf(":%d", port)
	//If you have more config values. Make sure they are not empty. If they are,
	//panic if your service requires the value to be non empty. There is no way
	//for the service to recover from it.
}
