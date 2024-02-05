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

build: build-dependencies
	tailwindcss -i src/styles.css -o public/styles.css
	templ generate
	go build -o redmage
