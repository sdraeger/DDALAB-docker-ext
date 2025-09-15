# Multi-stage build for Docker Desktop Extension

FROM golang:1.21-alpine AS backend-builder
WORKDIR /backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ddalab-manager .

FROM alpine:latest
RUN apk --no-cache add ca-certificates docker-cli docker-compose bash
COPY --from=backend-builder /backend/ddalab-manager /backend/ddalab-manager
COPY ui /ui
COPY metadata.json /metadata.json
COPY extension.yaml /extension.yaml
COPY icon.svg /icon.svg
COPY backend-entrypoint.sh /backend-entrypoint.sh

RUN chmod +x /backend/ddalab-manager /backend-entrypoint.sh

WORKDIR /

CMD ["/backend-entrypoint.sh"]

LABEL org.opencontainers.image.title="DDALAB Manager" \
      org.opencontainers.image.description="Manage DDALAB installations" \
      org.opencontainers.image.vendor="DDALAB Team" \
      com.docker.desktop.extension.api.version="0.3.4" \
      com.docker.desktop.extension.icon="https://raw.githubusercontent.com/sdraeger/DDALAB/main/logo.svg" \
      com.docker.extension.categories="utility,data-science" \
      com.docker.extension.screenshots="" \
      com.docker.extension.detailed-description="DDALAB Manager helps you manage DDALAB (Delay Differential Analysis Laboratory) installations. Monitor services, view logs, create backups, and control your DDALAB deployment directly from Docker Desktop." \
      com.docker.extension.publisher-url="https://github.com/sdraeger/DDALAB" \
      com.docker.extension.additional-urls="" \
      com.docker.extension.changelog=""