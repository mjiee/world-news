# app dependency check
check_dependency:
	@command -v wails >/dev/null 2>&1 || { \
		echo "Error: wails is not installed. Installing..."; \
		go install github.com/wailsapp/wails/v2/cmd/wails@latest; \
	}


# run app
run-app: check_dependency
	@wails dev

# generate db
db:
	@go run backend/repository/generator/main.go