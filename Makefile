IMAGE_NAME := raver
PLATFORM   := linux/arm/v7
SSH_HOST   := pi
ENV_FILE   := raver_env
CONTAINER_NAME := raver
LOCAL_PORT := 80
CONTAINER_PORT := 3000

.PHONY: run build save deploy all

run:
	go tool templ generate --watch --proxy="http://localhost:3000" --cmd="go run ."

build:
	docker buildx build --platform $(PLATFORM) -t $(IMAGE_NAME) --load .

deploy: build
	docker save $(IMAGE_NAME):latest | ssh $(SSH_HOST) "\
		docker load && \
		docker rm -f $(CONTAINER_NAME) 2>/dev/null || true && \
		docker run --env-file $(ENV_FILE) -d --name $(CONTAINER_NAME) -p $(LOCAL_PORT):$(CONTAINER_PORT) $(IMAGE_NAME):latest\
	"
