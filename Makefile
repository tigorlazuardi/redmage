.ONESHELL:

# Variables
export PATH := $(shell pwd)/node_modules/.bin:$(shell pwd)/bin:$(PATH)
export GOBIN := $(shell pwd)/bin
export GOOSE_DRIVER=sqlite3
export GOOSE_DBSTRING=./data.db
export GOOSE_MIGRATION_DIR=db/migrations

start: dev-dependencies
	@air

dev-dependencies: build-dependencies
	@if ! command -v air > /dev/null; then
		mkdir -p bin
		echo "Modd not found in PATH, installing it to $(shell pwd)/bin/air"
		go install github.com/cosmtrek/air@latest
	fi

build-dependencies:
	@if ! command -v templ > /dev/null; then
		mkdir -p bin
		echo "Templ not found in PATH, installing it to $(shell pwd)/bin/templ"
		go install github.com/a-h/templ/cmd/templ@v0.2.648
	fi
	@if [ ! -d "node_modules" ]; then
		echo "Node modules not found, installing them"
		npm install
	fi
	@if [  ! -f "public/htmx-1.9.11.min.js" ]; then
		mkdir -p public
		echo "Htmx not found, installing it"
		curl -o public/htmx-1.9.11.min.js https://unpkg.com/htmx.org@1.9.11/dist/htmx.min.js
	fi
	@if [ ! -f "public/htmx-response-targets-1.9.11.min.js" ]; then
		mkdir -p public
		echo "Htmx response targets not found, installing it"
		curl -o public/htmx-response-targets-1.9.11.min.js https://cdnjs.cloudflare.com/ajax/libs/htmx/1.9.11/ext/response-targets.min.js
	fi
	@if [ ! -f "public/dayjs-1.11.10.min.js" ]; then
		mkdir -p public
		echo "Dayjs not found, installing it"
		curl -o public/dayjs-1.11.10.min.js https://cdnjs.cloudflare.com/ajax/libs/dayjs/1.11.10/dayjs.min.js
	fi
	@if [ ! -f "public/dayjs-relativeTime-1.11.10.min.js" ]; then
		mkdir -p public
		echo "Dayjs Relative Time not found, installing it"
		curl -o public/dayjs-relativeTime-1.11.10.min.js https://cdnjs.cloudflare.com/ajax/libs/dayjs/1.11.10/plugin/relativeTime.min.js
	fi
	@if [ ! -f "public/theme-change-2.0.2.min.js" ];  then
		mkdir -p public
		echo "Theme change not found, installing it"
		curl -o public/theme-change-2.0.2.min.js https://cdn.jsdelivr.net/npm/theme-change@2.0.2/index.js
	fi

build: build-dependencies prepare
	go build -o redmage

prepare: gen
	mkdir -p public
	tailwindcss -i views/style.css -o public/style.css
	templ generate

gen:
	@go run github.com/stephenafamo/bob/gen/bobgen-sqlite@latest

migrate-new:
	@read -p "Name new migration: " name
	if [[ $$name ]]; then goose create "$$name" sql; fi
	
migrate-redo:
	@goose redo

migrate-up:
	@goose up
