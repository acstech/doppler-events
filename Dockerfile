# start the build process
FROM golang:latest as builder
WORKDIR /go/src/github.com/acstech/doppler-events/
<<<<<<< HEAD
<<<<<<< HEAD
# copy all the needed files for the program
COPY . .
ENV CGO_ENABLED=0
RUN go build -o ./grpcTest -ldflags "-s -w" github.com/acstech/doppler-events/cmd/grpcTEST
<<<<<<< HEAD
# move the build file into the final docker image
FROM alpine:latest
COPY --from=builder /go/src/github.com/acstech/doppler-events/grpcTest /opt/service/
COPY ./entrypoint.sh .
EXPOSE 8080
CMD ["./entrypoint.sh"] 
=======
# copy the rpc file to the build directory 
=======
# copy all the needed files for the program
>>>>>>> small deployment image
COPY . /go/src/github.com/acstech/doppler-events/
RUN CGO_ENABLED=0 GOOS=linux go build /go/src/github.com/acstech/doppler-events/cmd/grpcTEST/serviceStart.go
# move the build file into the final docker image
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/acstech/doppler-events/serviceStart /go/src/github.com/acstech/doppler-events/serviceStart
# copy dependencies and their trackers to the build directory
<<<<<<< HEAD
COPY vendor/* ./vendor/
COPY Gopkg.lock .
COPY Gopkg.toml .
CMD ["./backendServer"] 
>>>>>>> set up a starting point for the dockerfile
=======
COPY vendor/* /go/src/github.com/acstech/doppler-events/vendor/
COPY Gopkg.lock /go/src/github.com/acstech/doppler-events/Gopkg.lock
COPY Gopkg.toml /go/src/github.com/acstech/doppler-events/Gopkg.toml 
EXPOSE 8080
CMD ["/go/src/github.com/acstech/doppler-events/serviceStart"] 
>>>>>>> first working and very inefficient prototype
=======
# move the build file into the final docker image
FROM alpine:latest
COPY --from=builder /go/src/github.com/acstech/doppler-events/grpcTest /opt/service/
EXPOSE 8080
CMD ["//opt/service/grpcTest"] 
>>>>>>> Updated docker images
