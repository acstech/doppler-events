# start the build process
FROM golang:latest as builder
WORKDIR /go/src/github.com/acstech/doppler-events/
# copy all the needed files for the program
COPY . /go/src/github.com/acstech/doppler-events/
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