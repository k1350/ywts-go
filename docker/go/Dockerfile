FROM golang:1.14.2-alpine

ENV GO111MODULE=on

RUN apk update && \
    apk add curl git && \
    go get github.com/pilu/fresh