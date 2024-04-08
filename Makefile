PROJ_DIR 							?= $(shell pwd)
BIN_DIR               ?= $(PROJ_DIR)/bin
BINARY 								:= buildkit-exporter
BUILDX_BUILDER 				?= local
BUILDX_BUILDER_FILE 	?= ~/.docker/buildx/instances/$(BUILDX_BUILDER)
BUILDKIT_SOCKET_PATH 	?= /run/buildkit/buildkitd.sock
BUILDKIT_ADDR 				?= unix://$(BUILDKIT_SOCKET_PATH)
BUILDX_REMOTE_ADDR 		?= $(BUILDKIT_ADDR)
EXPORTER_ARGS 				 = --buildkit.address=$(BUILDKIT_ADDR) --web.listen-address=$(EXPORTER_ADDR)
EXPORTER_ADDR 				:= :9220
PROMU 								:= $(GOPATH)/bin/promu
PROMU_VERSION 				?= 0.15.0

IMAGE_NAME = buildkit-exporter

all: build run

$(BUILDX_BUILDER_FILE):
	docker buildx create \
		--name=$(BUILDX_BUILDER) \
		--use \
		--bootstrap \
		--driver=remote \
		$(BUILDX_REMOTE_ADDR)

build: $(PROMU)
	$(PROMU) build --prefix $(BIN_DIR)

run:
	$(BIN_DIR)/$(BINARY) $(EXPORTER_ARGS)

docker-build: $(BUILDX_BUILDER_FILE)
	docker buildx build \
		--builder=$(BUILDX_BUILDER) \
		--build-arg PROMU_VERSION=$(PROMU_VERSION) \
		--load \
		-t $(IMAGE_NAME) \
		.

.PHONY: run
docker-run: docker-build
	echo "ARGS: $(EXPORTER_ARGS)"
	docker run \
		--rm \
		-it \
		-v $(BUILDKIT_SOCKET_PATH):$(BUILDKIT_SOCKET_PATH) \
		-e EXPORTER_ARGS=$(EXPORTER_ARGS) \
		$(IMAGE_NAME)

$(PROMU):
	go install github.com/prometheus/promu@v$(PROMU_VERSION)

install-promu: $(PROMU)

run-prometheus:
	docker run \
		--network host \
    -v $(PROJ_DIR)/prometheus.yml:/etc/prometheus/prometheus.yml \
    prom/prometheus

get-histories:
	@buildctl --addr $(BUILDKIT_ADDR) debug histories --format '{{json .}}' \
		| jq -c '.' \
		| jq -s '.[0].record '
