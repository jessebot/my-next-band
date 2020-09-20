FROM golang:1.15

run mkdir -p /app

WORKDIR /app

ADD . /app

RUN go build ./bands.go
