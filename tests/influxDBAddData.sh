#!/bin/bash
for i in {1..1000}
do
    echo $i
    j=$(($i%100))
    docker-compose exec influx_db \
        curl -i -XPOST 'http://influx_db:8086/write?db=doppler' -d "influxTest,name=Bob,region=$j value=$j"
done
for i in {1..1000}
do
    echo $i
    j=$(($i%100))
    docker-compose exec influx_db \
        curl -i -XPOST 'http://influx_db:8086/write?db=doppler' -d "influxTest,name=Sally,region=$j value=$j"
done