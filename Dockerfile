FROM golang:1.11.10 AS build-env

WORKDIR /app

COPY . /app
RUN go build .
CMD ./auto-pocketer
