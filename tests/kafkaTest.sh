#!/bin/bash
echo "Kafka test begins..."
docker-compose exec kafka  \
  kafka-topics --create --topic "influx-topic" --partitions 1 --replication-factor 1 --if-not-exists --zookeeper zookeeper:2181
docker-compose exec kafka  \
  kafka-topics --describe --topic "influx-topic" --zookeeper zookeeper:2181
docker-compose exec kafka  \
  bash -c "echo '{\"clientID\":\"Bob\",\"eventID\":\"Sign In\",\"lat\":\"87687\",\"lon\":\"665.67676\",\"timeSinceEpoch\":\"2018-06-13T20:59:42.279355071Z\"}' | kafka-console-producer --request-required-acks 1 --broker-list kafka:9092 --topic 'influx-topic'"
docker-compose exec kafka  \
  bash -c "echo '{\"clientID\":\"Sally\",\"eventID\":\"Sign In\",\"lat\":\"87687\",\"lon\":\"665.67676\",\"timeSinceEpoch\":\"2018-06-13T20:59:42.279355071Z\"}' | kafka-console-producer --request-required-acks 1 --broker-list kafka:9092 --topic 'influx-topic'"
docker-compose exec kafka  \
  kafka-console-consumer --bootstrap-server kafka:9092 --topic "influx-topic" --from-beginning --max-messages 2
