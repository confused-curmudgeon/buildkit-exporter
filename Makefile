# Copyright 2024 Cody Boggs
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

# http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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
EXPORTER_VERSION 			:= $(shell cat VERSION)
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
		--rm \
		--network host \
    -v $(PROJ_DIR)/prometheus.yml:/etc/prometheus/prometheus.yml \
    prom/prometheus

get-histories:
	@buildctl --addr $(BUILDKIT_ADDR) debug histories --format '{{json .}}' \
		| jq -c '.' \
		| jq -s '.[0].record '

push-ghcr: docker-build
	@echo $(GITHUB_TOKEN_BUILDKIT_EXPORTER_GHCR) | docker login ghcr.io -u cboggs --password-stdin
	docker tag buildkit-exporter:latest ghcr.io/confused-curmudgeon/buildkit-exporter:v$(EXPORTER_VERSION)
	docker push ghcr.io/confused-curmudgeon/buildkit-exporter:v$(EXPORTER_VERSION)
