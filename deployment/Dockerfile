FROM golang:1.18 as modcache
WORKDIR /modcache
COPY go.mod go.sum /modcache/
RUN go mod download

FROM golang:1.18 AS builder
WORKDIR /app
COPY --from=modcache /go/pkg /go/pkg
COPY . /app
RUN cd /app && make build

FROM alpine:3.16.2
WORKDIR /app
COPY --from=builder /app/bin/antibruteforce /app/antibruteforce
CMD ["/app/antibruteforce"]