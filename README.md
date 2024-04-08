# Buildkit Exporter
Prometheus exporter for buildkitd metrics.

Supported buildkitd versions:
  - v0.12.4

## Building and running

### Build
```
# Install prereqs, build native binary, and run it with defaults
make

# Build native binary
make build

# Build docker image
make docker-build
```

### Run
```
# Run native binary
bin/buildkit-exporter --help

# Run docker image
make docker-run EXPORTER_ARGS=--help
```

## Flags
```
  -h, --[no-]help                Show context-sensitive help (also try --help-long and --help-man).
      --web.telemetry-path="/metrics"
                                 Path under which to expose metrics.
      --buildkit.address="unix:///run/buildkit/buildkitd.sock"
                                 Address to use for connecting to Buildkit
      --[no-]tls.insecure-skip-verify
                                 Ignore certificate and server verification when using a tls connection.
      --web.listen-address=:9220 ...
                                 Addresses on which to expose metrics and web interface. Repeatable for multiple addresses.
      --log.level=info           Only log messages with the given severity or above. One of: [debug, info, warn, error]
      --log.format=logfmt        Output format of log messages. One of: [logfmt, json]
      --[no-]version             Show application version.
```
