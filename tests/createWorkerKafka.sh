#!/bin/bash
docker-compose exec influx_db \
    curl -X POST -H "Content-Type: application/json" --data '{"name": "InfluxSink", "config": {"connector.class":"com.datamountaineer.streamreactor.connect.influx.InfluxSinkConnector", "tasks.max":"1", "topics":"influx-topic", "connect.influx.url":"http://influx_db:8086", "connect.influx.db":"mydb", "connect.influx.kcql":"INSERT INTO test SELECT * FROM influx-topic WITHTIMESTAMP sys_time()", "connect.influx.username": "root" }}' http://kafka_connect:8083/connectors | jq

docker-compose exec kafka \ bin/kafka-avro-console-producer \
  --broker-list kafka:9092 --topic influx-topic \
  --property value.schema='{"type":"record","name":"User",
  "fields":[{"name":"company","type":"string"},{"name":"address","type":"string"},{"name":"latitude","type":"float"},{"name":"longitude","type":"float"}]}'