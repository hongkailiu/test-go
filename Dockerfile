FROM docker.io/library/golang:1.24 as builder
WORKDIR /app

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /graph-builder ./cmd

FROM registry.access.redhat.com/ubi9/ubi:latest
COPY --from=builder /graph-builder /usr/bin/

ENTRYPOINT ["/usr/bin/graph-builder"]