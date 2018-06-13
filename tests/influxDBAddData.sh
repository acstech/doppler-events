#!/bin/bash
for i in {1..2}
do
    echo $i
    j=$(($i%10))
    docker-compose exec influx_db \
        curl -i -XPOST 'http://influx_db:8086/write?db=doppler' -d "influxTest,name=server$j,region=us-west value=$j"
done