# Doppler-Events
An API to connect clients to a queue to process for live preview.
## Setup
### Genral
- Clone the repository
- Install [go](https://golang.org/dl/) and [docker](https://docs.docker.com/install/)
- Run `docker-compose up -d`
- Map "kafka" to localhost in your host file
### Couchbase
- Go to localhost:8091
- Create an administrator
- Choose to configure the memory and choose 256 MB 
- Create a bucket
- Create a user that only has read and write permissions on that same bucket using the security tab
- Next copy and rename .env.defualt from data/couchbase/ to the base directory of this repository as .env.
- After that, fill in the appropriate environment varialbe.
## Testing **** Not completed yet ****
