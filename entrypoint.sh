#!/bin/sh
echo "Making sure kafka and couchbase are up"
while ! nc -z couchbase 8091 || ! nc -z kafka 9092 ;
do
    echo sleeping;
    sleep 1;
done;
echo "Kafka and Couchbase are up"
# safe to run the service
/opt/service/grpcTest
    