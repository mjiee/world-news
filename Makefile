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
module:
	@wails generate module

# build web static
.PHONY: build-web-static
build-web-static:
	@cd frontend && npm run build-web

# run app
.PHONY: run
run: check-dependency
	@wails dev

# run web
.PHONY: run-web
run-web: build-web-static
	@go run -tags web main.go

