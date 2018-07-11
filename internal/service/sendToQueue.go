package service

import (
	"fmt"
	"log"

	"github.com/Shopify/sarama"
)

//sendToQueue takes byte array, passes it to producer and writes to kafka instance
func (prod *Service) sendToQueue(JSONob []byte) {
	// create message to be sent to kafka
	log.Println(prod.kafkaTopic)
	msg := &sarama.ProducerMessage{
		Topic: prod.kafkaTopic,
		Value: sarama.ByteEncoder(JSONob),
		//Partition: 1, //partNum,
	}
	go func() {
		for err := range prod.producer.Errors() {
			fmt.Println("Failed to produce message:", err)
		}
	}()
	prod.producer.Input() <- msg
}
