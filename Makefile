.ONESHELL:

# Variables
export PATH := $(shell pwd)/node_modules/.bin:$(shell pwd)/bin:$(PATH)
export GOBIN := $(shell pwd)/bin
export GOOSE_DRIVER=sqlite3
export GOOSE_DBSTRING=./data.db
export GOOSE_MIGRATION_DIR=db/migrations

export REDMAGE_WEB_DEPENDENCIES_HTMX_VERSION=$(shell echo "$${REDMAGE_WEB_DEPENDENCIES_HTMX_VERSION:-1.9.12}")
export REDMAGE_WEB_DEPENDENCIES_DAYJS_VERSION=$(shell echo "$${REDMAGE_WEB_DEPENDENCIES_DAYJS_VERSION:-1.11.10}")
export REDMAGE_WEB_DEPENDENCIES_THEMECHANGE_VERSION=$(shell echo "$${REDMAGE_WEB_DEPENDENCIES_THEMECHANGE_VERSION:-2.0.2}")

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
	@if [  ! -f "public/htmx-${REDMAGE_WEB_DEPENDENCIES_HTMX_VERSION}.min.js" ]; then
		mkdir -p public
		echo "Htmx ${REDMAGE_WEB_DEPENDENCIES_HTMX_VERSION} not found, installing it"
		curl -o public/htmx-${REDMAGE_WEB_DEPENDENCIES_HTMX_VERSION}.min.js https://unpkg.com/htmx.org@${REDMAGE_WEB_DEPENDENCIES_HTMX_VERSION}/dist/htmx.min.js
	fi
	@if [ ! -f "public/htmx-response-targets-${REDMAGE_WEB_DEPENDENCIES_HTMX_VERSION}.min.js" ]; then
		mkdir -p public
		echo "Htmx response targets ${REDMAGE_WEB_DEPENDENCIES_HTMX_VERSION} not found, installing it"
		curl -o public/htmx-response-targets-${REDMAGE_WEB_DEPENDENCIES_HTMX_VERSION}.min.js https://cdnjs.cloudflare.com/ajax/libs/htmx/${REDMAGE_WEB_DEPENDENCIES_HTMX_VERSION}/ext/response-targets.min.js
	fi
	@if [ ! -f "public/dayjs-${REDMAGE_WEB_DEPENDENCIES_DAYJS_VERSION}.min.js" ]; then
		mkdir -p public
		echo "Dayjs ${REDMAGE_WEB_DEPENDENCIES_DAYJS_VERSION} not found, installing it"
		curl -o public/dayjs-${REDMAGE_WEB_DEPENDENCIES_DAYJS_VERSION}.min.js https://cdnjs.cloudflare.com/ajax/libs/dayjs/${REDMAGE_WEB_DEPENDENCIES_DAYJS_VERSION}/dayjs.min.js
	fi
	@if [ ! -f "public/dayjs-relativeTime-${REDMAGE_WEB_DEPENDENCIES_DAYJS_VERSION}.min.js" ]; then
		mkdir -p public
		echo "Dayjs Relative Time ${REDMAGE_WEB_DEPENDENCIES_DAYJS_VERSION} not found, installing it"
		curl -o public/dayjs-relativeTime-${REDMAGE_WEB_DEPENDENCIES_DAYJS_VERSION}.min.js https://cdnjs.cloudflare.com/ajax/libs/dayjs/${REDMAGE_WEB_DEPENDENCIES_DAYJS_VERSION}/plugin/relativeTime.min.js
	fi
	@if [ ! -f "public/dayjs-utc-${REDMAGE_WEB_DEPENDENCIES_DAYJS_VERSION}.min.js" ]; then
		mkdir -p public
		echo "Dayjs UTC ${REDMAGE_WEB_DEPENDENCIES_DAYJS_VERSION} plugin not found, installing it"
		curl -o public/dayjs-utc-${REDMAGE_WEB_DEPENDENCIES_DAYJS_VERSION}.min.js https://cdnjs.cloudflare.com/ajax/libs/dayjs/${REDMAGE_WEB_DEPENDENCIES_DAYJS_VERSION}/plugin/utc.min.js
	fi
	@if [ ! -f "public/dayjs-timezone-${REDMAGE_WEB_DEPENDENCIES_DAYJS_VERSION}.min.js" ]; then
		mkdir -p public
		echo "Dayjs Timezone ${REDMAGE_WEB_DEPENDENCIES_DAYJS_VERSION} plugin not found, installing it"
		curl -o public/dayjs-timezone-${REDMAGE_WEB_DEPENDENCIES_DAYJS_VERSION}.min.js https://cdnjs.cloudflare.com/ajax/libs/dayjs/${REDMAGE_WEB_DEPENDENCIES_DAYJS_VERSION}/plugin/timezone.min.js
	fi
	@if [ ! -f "public/theme-change-${REDMAGE_WEB_DEPENDENCIES_THEMECHANGE_VERSION}.min.js" ];  then
		mkdir -p public
		echo "Theme change ${REDMAGE_WEB_DEPENDENCIES_THEMECHANGE_VERSION} not found, installing it"
		curl -o public/theme-change-${REDMAGE_WEB_DEPENDENCIES_THEMECHANGE_VERSION}.min.js https://cdn.jsdelivr.net/npm/theme-change@${REDMAGE_WEB_DEPENDENCIES_THEMECHANGE_VERSION}/index.js
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
	
migrate-down:
	@goose down

migrate-up:
	@goose up
