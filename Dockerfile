# Build Geth in a stock Go builder container
FROM golang:1.11-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers

ADD . /go/src/Smilo-blackbox
RUN cd /go/src/Smilo-blackbox && make build

# Pull Geth into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/Smilo-blackbox/blackbox /usr/local/bin/

EXPOSE 9000
ENTRYPOINT ["blackbox"]
