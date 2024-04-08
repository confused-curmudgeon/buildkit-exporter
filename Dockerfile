FROM golang:1.21.6-alpine AS builder
ADD . /build
WORKDIR /build

RUN apk add git make
RUN --mount=type=cache,dst=/go/bin \
    make install-promu

RUN --mount=type=cache,dst=/go/pkg/mod \
    --mount=type=cache,dst=/root/.cache/go-build \
    make build

FROM alpine
LABEL org.opencontainers.image.source https://github.com/confused-curmudgeon/buildkit-exporter
COPY --from=builder /build/bin/buildkit-exporter /buildkit-exporter
ENTRYPOINT /buildkit-exporter ${EXPORTER_ARGS}
