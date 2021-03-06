FROM golang:1.14.2

ENV \
  GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

WORKDIR /go/src/github.com/egeneralov/dockerd-health-exporter
ADD go.mod go.sum /go/src/github.com/egeneralov/dockerd-health-exporter/
RUN go mod download

ADD . .

RUN go build -v -installsuffix cgo -ldflags="-w -s" -o /go/bin/dockerd-health-exporter .


FROM debian:buster

ENV PATH='/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin'
CMD /go/bin/dockerd-health-exporter

COPY --from=0 /go/bin /go/bin
