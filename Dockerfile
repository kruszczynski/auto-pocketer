FROM golang:1.11.5

WORKDIR /go/src/github.com/kruszczynski/auto-pocketer

COPY . /go/src/github.com/kruszczynski/auto-pocketer
RUN go get -u
