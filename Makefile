localBuilder = ~/.docker/buildx/instances/local
IMAGE_NAME = buildkit-exporter

all: build run

$(localBuilder):
	docker buildx create \
		--name=local \
		--use \
		--bootstrap \
		--driver=remote \
		unix://${XDG_RUNTIME_DIR}/buildkit.sock

build: $(localBuilder)
	docker buildx build \
		--builder=local \
		--load \
		-t $(IMAGE_NAME) \
		.

run:
	docker run -it \
		-v $(XDG_RUNTIME_DIR)/buildkit.sock:/run/buildkit.sock \
		--user 1000 \
		$(IMAGE_NAME)
