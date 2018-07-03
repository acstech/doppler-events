#! /bin/bash
echo "Start up couchbase and kafka..."
docker-compose run wait
echo "Startup the rest of the services..."
docker-compose up -d
echo "Start up complete!"