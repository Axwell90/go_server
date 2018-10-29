FROM golang:1.9 as builder

RUN mkdir /go/src/go_server
WORKDIR /go/src/go_server

COPY . .

RUN go get -v
RUN go build -v

FROM ubuntu:16.04

ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /root/
COPY . .

COPY --from=builder /go/src/go_server/go_server ./go_server

RUN chmod 755 go_server

ENTRYPOINT ["/root/go_server"]