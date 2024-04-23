FROM golang:1.19.1-alpine

RUN apk update && apk add git

RUN mkdir /go/src/github.com/backend

WORKDIR /go/src/github.com/backend

ADD . /go/src/github.com/backend

RUN go get -u golang.org/x/tools/cmd/goimports
