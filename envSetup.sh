#!/bin/bash

#Run this in doppler-events with the other doppler folders in the same
#directory as doppler-events

#saves current directory for future reference in other terminal windows
cwd=`pwd`

#start docker
docker-compose up -d

#let docker-compose warm up and get services running
wait
sleep 45

# run front-end api server
osascript -e 'tell application "Terminal" to do script "cd '$cwd' && cd ../doppler-api/ && go run cmd/doppler-api/main.go"'

# serve front end to localhost:9080 (or whatever your static-server is set to serve to)
osascript -e 'tell application "Terminal" to do script "cd '$cwd' && cd ../doppler-frontend/ && static-server"'

#run backend API (producer)
go run cmd/grpcTEST/serviceStart.go

