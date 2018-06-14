#!/bin/bash
echo "Kafka connect connectors..."
docker-compose exec kafka_connect  \
    curl kafka_connect:8083/connector-plugins | jq
echo "Kafka plugins"
docker-compose exec kafka_connect  \
    curl kafka_connect:8083/connector-plugins | jq 
echo "Kafka active connectors"
docker-compose exec kafka_connect  \
    curl kafka_connect:8083/connectors
echo "Kafka connector status for Influx Connector"
docker-compose exec kafka_connect  \
    curl kafka_connect:8083/connectors/InfluxSink/status | jq