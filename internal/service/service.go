// Event API Server

package service

import (
	"github.com/Shopify/sarama"
	cb "github.com/acstech/doppler-events/internal/couchbase"
)

// Struct to hold parameters for server
type Service struct {
	errorResponse ErrorRes
	producer      sarama.AsyncProducer
	CbConn        *cb.Couchbase
	kafkaTopic    string
}

// Error struct for verifying constraints
type ErrorRes struct {
	errMes []string
}

// Init sets up the backend server so that clients can send data to kafka
// cbCon is the connection string for couchbase
// returns an error if any occur while creating a kafka producer, a couchbase connection, sending data,
// or closing the kafka producer
func NewService(producer sarama.AsyncProducer, cbConn *cb.Couchbase, kafkaTopic string) (*Service, error) {
	// Create instance of empty ErrorRes struct
	e := ErrorRes{}

	return &Service{
		errorResponse: e,
		producer:      producer,
		CbConn:        cbConn,
		kafkaTopic:    kafkaTopic,
	}, nil
}
