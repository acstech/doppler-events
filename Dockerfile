# start the build process
FROM golang:latest as builder
WORKDIR /go/src/github.com/acstech/doppler-events/
# copy the rpc file to the build directory 
COPY . /go/src/github.com/acstech/doppler-events/
# copy the backend server starter to the build directory
# COPY cmd/grpcTEST/serviceStart.go /go/src/github.com/acstech/doppler-events/cmd/grpcTEST/serviceStart.go
# # # copy the couchbase connector to the build directory
# COPY internal/couchbase/couchbaseConn.go /go/src/github.com/acstech/doppler-events/internal/couchbase/couchbaseConn.go
# # # copy the backend server to the build directory
# COPY internal/service/service.go /go/src/github.com/acstech/doppler-events/internal/service/service.go
# # copy dependencies and their trackers to the build directory
# COPY vendor/* ./vendor/
# COPY Gopkg.lock .
# COPY Gopkg.toml .
RUN CGO_ENABLED=0 GOOS=linux go build /go/src/github.com/acstech/doppler-events/cmd/grpcTEST/serviceStart.go
# move the build file into the final docker image
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/acstech/doppler-events/serviceStart /go/src/github.com/acstech/doppler-events/serviceStart
# copy dependencies and their trackers to the build directory
COPY vendor/* /go/src/github.com/acstech/doppler-events/vendor/
COPY Gopkg.lock /go/src/github.com/acstech/doppler-events/Gopkg.lock
COPY Gopkg.toml /go/src/github.com/acstech/doppler-events/Gopkg.toml 
EXPOSE 8080
CMD ["/go/src/github.com/acstech/doppler-events/serviceStart"] 