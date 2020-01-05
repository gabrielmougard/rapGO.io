FROM golang:latest

WORKDIR $GOPATH/src/rapGO.io/src/audioconverterservice
COPY . $GOPATH/src/rapGO.io/src/audioconverterservice
RUN go build .

EXPOSE 3001
ENTRYPOINT ["./audioconverterservice"]