APP_NAME=world-news

# app dependency check
check-dependency:
	@command -v wails >/dev/null 2>&1 || { \
		echo "Error: wails is not installed. Installing..."; \
		go install github.com/wailsapp/wails/v2/cmd/wails@latest; \
	}

# generate db
.PHONY: db
db:
	@go run backend/repository/generator/main.go

# manually generate the wailsjs directory
.PHONY: module
module: check-dependency
	@wails generate module

# build web static
.PHONY: build-web-static
build-web-static:
	@cd frontend && npm run build-web

# run dev
.PHONY: dev
dev: check-dependency
	@wails dev

# run web
.PHONY: run-web
run-web: build-web-static
	@go run -tags web main.go


# build app
.PHONY: build
build: check-dependency
	@wails build -clean -ldflags "-s -w"

# build web
.PHONY: build-web
build-web: build-web-static
	@go build -tags web -o ./build/bin/$(BIN) main.go

# deploy
.PHONY: deploy
deploy:
	@docker compose up -d
