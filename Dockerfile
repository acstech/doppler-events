# start the build process
FROM golang:1.10 as builder
WORKDIR /go/src/github.com/acstech/doppler-events/
# copy all the needed files for the program
COPY . .
ENV CGO_ENABLED=0
RUN go build -o ./grpcTest -ldflags "-s -w" github.com/acstech/doppler-events/cmd/grpcTEST
# move the build file into the final docker image
FROM alpine:3.8
COPY --from=builder /go/src/github.com/acstech/doppler-events/grpcTest /opt/service/
COPY ./entrypoint.sh .
EXPOSE 8080
<<<<<<< HEAD
CMD ["./entrypoint.sh"] 
=======
CMD ["./entrypoint.sh"] 
>>>>>>> Updated docker files and made sure that they all work together
