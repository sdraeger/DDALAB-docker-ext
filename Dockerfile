# Multi-stage build for Docker Desktop Extension

# Build backend
FROM golang:1.23-alpine AS backend-builder
WORKDIR /backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ddalab-manager .

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates docker-cli docker-compose bash wget

# Create necessary directories and set permissions
RUN mkdir -p /run/guest-services && chmod 755 /run/guest-services

COPY --from=backend-builder /backend/ddalab-manager /backend/ddalab-manager
COPY ui/dist /ui
COPY metadata.json /metadata.json
COPY extension.yaml /extension.yaml
COPY compose.yaml /compose.yaml
COPY icon.svg /icon.svg

RUN chmod +x /backend/ddalab-manager

# Expose port 8080 for the backend
EXPOSE 8080

WORKDIR /

# Start the backend directly
LABEL org.opencontainers.image.title="DDALABManager" \
      org.opencontainers.image.description="Manage DDALAB installations" \
      org.opencontainers.image.vendor="DDALAB Team" \
      com.docker.desktop.extension.api.version="0.3.4" \
      com.docker.desktop.extension.icon="data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjQiIGhlaWdodD0iMjQiIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHBhdGggZD0iTTEyIDJMMTMuMDkgOC4yNkwyMCA5TDEzLjA5IDE1Ljc0TDEyIDIyTDEwLjkxIDE1Ljc0TDQgOUwxMC45MSA4LjI2TDEyIDJaIiBmaWxsPSIjNDA5NkZGIi8+Cjwvc3ZnPgo=" \
      com.docker.extension.categories="utility,data-science" \
      com.docker.extension.screenshots="[]" \
      com.docker.extension.detailed-description="DDALAB Manager helps you manage DDALAB installations. Monitor services, view logs, create backups, and control your DDALAB deployment directly from Docker Desktop." \
      com.docker.extension.publisher-url="https://github.com/sdraeger/DDALAB" \
      com.docker.extension.additional-urls="[]" \
      com.docker.extension.changelog="Initial release"

CMD ["/backend/ddalab-manager"]