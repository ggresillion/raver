IMAGE_NAME := raver
PLATFORM   := linux/arm/v7
SSH_HOST   := pi
ENV_FILE   := raver_env
CONTAINER_NAME := raver
LOCAL_PORT := 80
CONTAINER_PORT := 3000

.PHONY: dev dev/templ dev/server dev/tailwind dev/sync_assets assets build deploy

dev/templ:
	go tool templ generate --watch --proxy="http://localhost:$(CONTAINER_PORT)" --open-browser=false -v

dev/server:
	go run github.com/air-verse/air@v1.63.0 \
		--build.cmd "go build -o tmp/bin/main ." --build.bin "tmp/bin/main" \
		--build.delay "100" --build.exclude_dir "web/node_modules" \
		--build.include_ext "go" --build.stop_on_error "false" \
		--misc.clean_on_exit true

dev/tailwind:
	npx --yes @tailwindcss/cli -i ./web/static/css/input.css -o ./web/static/css/app.css --minify --watch

dev/sync_assets:
	go run github.com/air-verse/air@v1.63.0 \
		--build.cmd "go tool templ generate --notify-proxy" --build.bin "/usr/bin/true" \
		--build.delay "100" --build.exclude_dir "" \
		--build.include_dir "web/static" --build.include_ext "js,css"

dev:
	make -j4 dev/templ dev/server dev/tailwind dev/sync_assets

assets:
	npx --yes @tailwindcss/cli -i ./web/static/css/input.css -o ./web/static/css/app.css --minify

build: assets
	docker buildx build --platform $(PLATFORM) -t $(IMAGE_NAME) --load .

deploy: build
	docker save $(IMAGE_NAME):latest | ssh $(SSH_HOST) "\
		docker load && \
		docker rm -f $(CONTAINER_NAME) 2>/dev/null || true && \
		docker run --env-file $(ENV_FILE) -d --name $(CONTAINER_NAME) -p $(LOCAL_PORT):$(CONTAINER_PORT) $(IMAGE_NAME):latest\
	"
