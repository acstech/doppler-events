package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/Shopify/sarama"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	// Setup configuration
	config := sarama.NewConfig()
	// Return specifies what channels will be populated.
	// If they are set to true, you must read from
	config.Producer.Return.Successes = true
	// The level of acknowledgement reliability needed from the broker.
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	brokers := []string{"localhost:9092"}
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		// Should not reach here
		panic(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			// Should not reach here
			panic(err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	themessage := &message{
		ClientID: "abcd",
		EventID:  "theIDevent",
		Lat:      "34.1954",
		Long:     "-79.7626",
	}

	var enqueued, errors int
	doneCh := make(chan struct{})
	go func() {
		for {

			time.Sleep(500 * time.Millisecond)

			themessage.TimeSinceEpoch = time.Now().Unix()
			buf := new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(themessage)
			if err != nil {
				panic(err)
			}
			spew.Dump(buf.String())
			msg := &sarama.ProducerMessage{
				Topic: "influx-topic",
				Value: sarama.ByteEncoder(buf.Bytes()),
			}
			select {
			case producer.Input() <- msg:
				enqueued++
				fmt.Println("Produce message")
			case err := <-producer.Errors():
				errors++
				fmt.Println("Failed to produce message:", err)
			case <-signals:
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	log.Printf("Enqueued: %d; errors: %d\n", enqueued, errors)
}

type message struct {
	Lat            string `json:"lat"`
	Long           string `json:"lon"`
	TimeSinceEpoch int64  `json:"timeSinceEpoch"`
	ClientID       string `json:"clientID"`
	EventID        string `json:"eventID"`
}
