FROM golang:1.21.6-alpine AS builder
ADD . /build
WORKDIR /build
RUN --mount=type=cache,dst=/go/pkg/mod \
    --mount=type=cache,dst=/root/.cache/go-build \
    go get ./... && \
    go mod tidy && \
    go build -o buildkit-exporter ./src

FROM alpine
COPY --from=builder /build/buildkit-exporter /buildkit-exporter
ENTRYPOINT /buildkit-exporter
