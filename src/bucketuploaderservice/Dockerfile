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
ENV PROJECT rapGO.io/src/bucketuploaderservice
WORKDIR /go/src/$PROJECT
ADD glide.yaml glide.yaml
ADD glide.lock glide.lock
# install dependencies
RUN glide install
# add the source code
COPY . .
RUN go build -o bucketuploader
RUN ls
# strip and compress the binary
RUN strip --strip-unneeded bucketuploader
RUN upx bucketuploader

# use a minimal alpine image
FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
# set working directory
WORKDIR /root
# copy some needed binaries for librdkafka
COPY --from=builder /usr/lib/pkgconfig /usr/lib/pkgconfig
COPY --from=builder /usr/lib/librdkafka* /usr/lib/
# copy the binary from builder
COPY --from=builder /go/src/rapGO.io/src/bucketuploaderservice/bucketuploader .
# copy the private key for accessing the google cloud storage
COPY --from=builder /go/src/rapGO.io/src/bucketuploaderservice/keen-dispatch.json .

# define all the needed environment variables
ENV INPUT_PREFIX input
ENV INPUT_SUFFIX .mp3
ENV OUTPUT_PREFIX output
ENV OUTPUT_SUFFIX .mp3

ENV TMP_FOLDER /data/tmp/

ENV KAFKA_BROKER 
ENV KAFKA_CONSUMERGROUP_ID
ENV KAFKA_TOPIC_TOBUCKET toBucket
ENV KAFKA_TOPIC_TOCORE toCore
ENV KAFKA_TOPICPREFIX_HEARTBEAT heartbeat_
ENV KAFKA_TOPIC_TOSTREAM toStream

ENV STORAGE_PROJECT_ID keen-dispatch-263413
ENV STORAGE_BUCKET_NAME rapgo-bucket-1
ENV GOOGLE_APPLICATION_CREDENTIALS /root/keen-dispatch.json

# run the binary
ENTRYPOINT ["./bucketuploader"]