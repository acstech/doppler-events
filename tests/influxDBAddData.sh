#!/bin/bash
for i in {1..100}
do
    echo $i"\n"
    j=$(($i%10))
    docker-compose exec influx_db \
        curl -i -XPOST 'http://localhost:8086/write?db=mydb' --data-binary "server,host=$j,region=us-west value=0.64"
done