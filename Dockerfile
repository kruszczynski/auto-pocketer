FROM golang:1.11.5 AS build-env

WORKDIR /app

COPY . /app
RUN go build .
CMD ./auto-pocketer
