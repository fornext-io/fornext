FROM golang:1.20.2-bullseye as builder
COPY . /go/src/github.com/fornext-io/fornext
WORKDIR /go/src/github.com/fornext-io/fornext
RUN make build-linux

FROM debian:bullseye
RUN mkdir -p /app && \
    chown -R nobody:nogroup /app
COPY --from=builder /go/src/github.com/fornext-io/fornext/bin/golang-project-template /app
USER nobody
WORKDIR /app
ENTRYPOINT ["/app/golang-project-template"]
