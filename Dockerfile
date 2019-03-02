FROM golang:1.11.5 AS build-env

WORKDIR /app

COPY go.mod go.sum /app/
RUN go get -u -v ./...

COPY . /app
RUN go build .
CMD ./auto-pocketer
