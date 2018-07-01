# start the build process
FROM golang:latest as builder
WORKDIR /go/src/github.com/acstech/doppler-events/
# copy the rpc file to the build directory 
COPY rpc/eventAPI/eventAPI.pb.go . 
# copy the backend server starter to the build directory
COPY cmd/grpcTEST/serviceStart.go .
# copy the couchbase connector to the build directory
COPY internal/couchbase/couchbaseConn.go .
# copy the backend server to the build directory
COPY internal/service/service.go .
# copy dependencies and their trackers to the build directory
COPY vendor/* ./vendor/
COPY Gopkg.lock .
COPY Gopkg.toml .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o backendServer .
# move the build file into the final docker image
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/acstech/doppler-events/backendServer .
# copy dependencies and their trackers to the build directory
COPY vendor/* ./vendor/
COPY Gopkg.lock .
COPY Gopkg.toml .
CMD ["./backendServer"] 