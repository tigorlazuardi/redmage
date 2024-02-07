.ONESHELL:

# Variables
export PATH := $(shell pwd)/node_modules/.bin:$(shell pwd)/bin:$(PATH)
export GOBIN := $(shell pwd)/bin

start: dev-dependencies
	@modd

dev-dependencies: build-dependencies
	@if ! command -v modd > /dev/null; then
		echo "Modd not found in PATH, installing it to $(shell pwd)/bin/modd"
		go install github.com/cortesi/modd@latest
	fi

build-dependencies:
	@if ! command -v templ > /dev/null; then
		echo "Templ not found in PATH, installing it to $(shell pwd)/bin/templ"
		go install github.com/a-h/templ/cmd/templ@latest
	fi
	@if [ ! -d "node_modules" ]; then
		echo "Node modules not found, installing them"
		npm install
	fi
	@if [  ! -f "public/htmx-1.9.10.min.js" ]; then
		echo "Htmx not found, installing it"
		curl -o public/htmx-1.9.10.min.js https://unpkg.com/htmx.org@1.9.10/dist/htmx.min.js
	fi
	@if [ ! -f "public/htmx-response-targets-1.9.10.min.js" ]; then
		echo "Htmx response targets not found, installing it"
		curl -o public/htmx-response-targets-1.9.10.min.js https://cdnjs.cloudflare.com/ajax/libs/htmx/1.9.10/ext/response-targets.min.js
	fi

build: build-dependencies
	tailwindcss -i src/styles.css -o public/styles.css
	templ generate
	go build -o redmage
