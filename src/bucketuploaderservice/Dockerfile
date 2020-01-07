FROM golang:1.12-alpine AS builder
RUN apk add --no-cache ca-certificates git && \
      wget -qO/go/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && \
      chmod +x /go/bin/dep

ENV PROJECT rapGO.io/src/bucketuploaderservice
WORKDIR /go/src/$PROJECT

#restore dependencies
COPY Gopkg.* ./
RUN dep ensure --vendor-only -v

COPY . .
RUN go build -o /bucketuploaderservice .

FROM alpine AS release
RUN apk add --no-cache ca-certificates
WORKDIR /bucketuploaderservice
COPY --from=builder /bucketuploaderservice ./build

ENTRYPOINT ["/bucketuploaderservice/build"]