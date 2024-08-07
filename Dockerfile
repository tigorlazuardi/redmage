FROM node:22-bullseye AS web-builder
WORKDIR /web
COPY package.json package-lock.json ./
RUN npm install
COPY Makefile ./
RUN make web-dependencies
COPY views ./views
COPY tailwind.config.js ./
RUN make web-build

FROM devopsworks/golang-upx:1.22.1 AS builder
WORKDIR /app
COPY Makefile ./
RUN make build-dependencies
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=web-builder /web/public ./public
ARG REDMAGE_RUNTIME_VERSION=unknown
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    make build-docker && \
    strip /app/redmage && \
    /usr/local/bin/upx -9 /app/redmage

FROM gcr.io/distroless/base:latest
WORKDIR /app
COPY --from=builder /app/redmage /app/redmage
ENV REDMAGE_FLAGS_CONTAINERIZED=true
ENV REDMAGE_DB_STRING=/app/db/data.db
ENV REDMAGE_PUBSUB_DB_NAME=/app/db/pubsub.db
ENV REDMAGE_DOWNLOAD_DIRECTORY=/app/downloads
ENV REDMAGE_RUNTIME_ENVIRONMENT=production
CMD [ "/app/redmage", "serve" ]
