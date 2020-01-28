FROM golang:1.12-buster AS builder
# install xz
RUN apt-get update && apt-get install -y \
    xz-utils \
&& rm -rf /var/lib/apt/lists/*
# install UPX
ADD https://github.com/upx/upx/releases/download/v3.94/upx-3.94-amd64_linux.tar.xz /usr/local
RUN xz -d -c /usr/local/upx-3.94-amd64_linux.tar.xz | \
    tar -xOf - upx-3.94-amd64_linux/upx > /bin/upx && \
    chmod a+x /bin/upx

# install glide
RUN go get github.com/Masterminds/glide
# setup the working directory
ENV PROJECT rapGO.io/src/heartbeatservice
WORKDIR /go/src/$PROJECT

### dependencies
COPY vendor /go/src

# in production, comment the line above and un-comment the lines below.
#ADD glide.yaml glide.yaml
#ADD glide.lock glide.lock
## install dependencies
#RUN glide install

###############

# add the source code
COPY . .
RUN rm -r vendor
RUN ls
RUN go build -o heartbeat
# strip and compress the binary
RUN strip --strip-unneeded heartbeat
RUN upx heartbeat

# use a minimal alpine image
#FROM alpine
#RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
##create the TMP_FOLDER
#RUN mkdir -p /data/tmp
## set working directory
#WORKDIR /root
## copy some needed binaries
#COPY --from=builder /usr/lib/pkgconfig /usr/lib/pkgconfig
## copy the binary from builder
#COPY --from=builder /go/src/rapGO.io/src/heartbeatservice/heartbeat .

# ENV variables definition
ENV KAFKA_BROKER kafka:9093
ENV KAFKA_TOPIC_TOHEARTBEAT toHeartbeat
ENV KAFKA_TOPIC_TOBUCKET toBucket

ENV LAST_HEARTBEAT_DESC Rap generated successfully !|Saving generated data to cloud...
ENV TOTAL_HEARTBEAT_NUMBER 6
#

# run the binary
ENTRYPOINT ["./heartbeat"]