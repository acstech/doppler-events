package service

import (
	"github.com/Shopify/sarama"
	"log"
	"os"
	"time"
)

//newProducer configures an asynchronous kafka producer client, returns it
func newProducer(address string) (sarama.AsyncProducer, error) {
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	// Setup configuration
	config := sarama.NewConfig()
	config.ClientID = "1"
	//configuration for batches
	config.Producer.Flush.MaxMessages = 30 // will flush if 30 messages arrived
	config.Producer.Flush.Frequency = 50 * time.Millisecond
	config.Producer.Flush.Messages = 1 // can flush with 1 message
	//The level of acknowledgement reliability needed from the broker.
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// The level of acknowledgement reliability needed from the broker.
	config.Producer.RequiredAcks = sarama.WaitForAll
	brokers := []string{address}
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		// Should not reach here
		return nil, err
	}
	return producer, nil
}
