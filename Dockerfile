FROM golang:1.14.0 AS build-env

WORKDIR /app

COPY . /app
RUN go get .
RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' .

FROM alpine:3.9.4
RUN apk update && apk add ca-certificates
WORKDIR /app
COPY --from=build-env /app/auto-pocketer .
CMD ["./auto-pocketer"]
