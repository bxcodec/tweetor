## Builder
FROM golang:1.10-alpine3.7 as builder
RUN apk update && apk upgrade && \
    apk --no-cache --update add git make && \
    go get -u github.com/golang/dep/cmd/dep
WORKDIR /go/src/github.com/bxcodec/tweetor
COPY . .
RUN make engine

## Distribution
FROM alpine:latest
RUN apk update && apk upgrade && \
    apk --no-cache --update add ca-certificates tzdata && \
    mkdir /tweetor && mkdir /app

WORKDIR /tweetor

EXPOSE 9090

COPY --from=builder /go/src/github.com/bxcodec/tweetor/engine /app

CMD /app/engine