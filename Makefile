BUILDX_BUILDER ?= local
BUILDX_BUILDER_FILE = ~/.docker/buildx/instances/$(BUILDX_BUILDER)
BUILDKIT_SOCKET_PATH ?= /run/buildkit/buildkitd.sock
BUILDKIT_ADDR ?= unix://$(BUILDKIT_SOCKET_PATH)
BUILDX_REMOTE_ADDR ?= $(BUILDKIT_ADDR)

IMAGE_NAME = buildkit-exporter

all: build run

$(BUILDX_BUILDER_FILE):
	docker buildx create \
		--name=$(BUILDX_BUILDER) \
		--use \
		--bootstrap \
		--driver=remote \
		$(BUILDX_REMOTE_ADDR)

build: $(BUILDX_BUILDER_FILE)
	docker buildx build \
		--builder=$(BUILDX_BUILDER) \
		--load \
		-t $(IMAGE_NAME) \
		.

.PHONY: run
run:
	$(eval EXP_ARGS="-buildkit-addr $(BUILDKIT_ADDR)")
	echo $(EXP_ARGS)
	docker run \
		--rm \
		-it \
		-v $(BUILDKIT_SOCKET_PATH):$(BUILDKIT_SOCKET_PATH) \
		-e EXP_ARGS=$(EXP_ARGS) \
		$(IMAGE_NAME)
