FROM golang:1.15

run mkdir -p /app

WORKDIR /app

ADD . /app

RUN go get -u github.com/gorilla/mux

RUN go build ./bands.go

EXPOSE 8080

CMD ["./bands"]
