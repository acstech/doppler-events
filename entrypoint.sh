#!/bin/sh
echo "Making sure kafka and couchbase are up"
while [ ! nc -z kafka 9092 ] && [ ! nc -z couchbase 8091 ];
do
    echo sleeping;
    sleep 1;
done;
# safe to run the service
/opt/service/grpcTest
    