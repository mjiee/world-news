APP_NAME=world-news
MINIAUDIO_FILE="backend/pkg/audio/miniaudio.h"

export CGO_ENABLED=1

# app dependency check and install
.PHONY: init
init:
	@command -v wails >/dev/null 2>&1 || { \
		echo "Error: wails is not installed. Installing..."; \
		go install github.com/wailsapp/wails/v2/cmd/wails@latest; \
	}
	@if [ ! -f $(MINIAUDIO_FILE) ]; then \
		echo "miniaudio.h not found, downloading to $(MINIAUDIO_FILE)..."; \
		curl -L -o $(MINIAUDIO_FILE) https://raw.githubusercontent.com/mackron/miniaudio/master/miniaudio.h; \
	fi

# generate db
.PHONY: db
db:
	@go run backend/repository/generator/main.go

# manually generate the wailsjs directory
.PHONY: module
module: init
	@wails generate module

# build web static
.PHONY: build-web-static
build-web-static:
	@cd frontend && npm run build-web

# run dev
.PHONY: dev
dev: init
	@wails dev

# run web
.PHONY: run-web
run-web: build-web-static
	@go run -tags web main.go


# build app
.PHONY: build
build: init
	@wails build -clean -ldflags "-s -w" 

# build web
.PHONY: build-web
build-web: build-web-static
	@go build -tags web -o ./build/bin/$(BIN) main.go

# deploy
.PHONY: deploy
deploy:
	@docker compose up -d
