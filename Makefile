# app dependency check
check-dependency:
	@command -v wails >/dev/null 2>&1 || { \
		echo "Error: wails is not installed. Installing..."; \
		go install github.com/wailsapp/wails/v2/cmd/wails@latest; \
	}


# run app
run-app: check-dependency
	@wails dev

# generate db
.PHONY: db
db:
	@go run backend/repository/generator/main.go