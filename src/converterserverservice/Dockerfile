FROM golang:latest

WORKDIR $GOPATH/src/rapGO.io/src/audioconverterservice
COPY . $GOPATH/src/rapGO.io/src/audioconverterservice
RUN go build .

EXPOSE 3001

# env setup
ENV SERVER_HTTP_PORT 3001
ENV APP_TMP_FOLDER  
#
ENTRYPOINT ["./converterserverservice"]