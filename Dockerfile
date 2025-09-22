# Multi-stage build for Docker Desktop Extension

# Use pre-built backend from ddalab-control
FROM sdraeger/ddalab-control:latest AS backend

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates docker-cli docker-compose bash wget

# Create necessary directories and set permissions
RUN mkdir -p /run/guest-services && chmod 755 /run/guest-services

COPY --from=backend /app/ddalab-control /backend/ddalab-manager
COPY ui/dist /ui
COPY metadata.json /metadata.json
COPY extension.yaml /extension.yaml
COPY extension-manifest.json /extension-manifest.json
COPY compose.yaml /compose.yaml
COPY icon.svg /icon.svg

RUN chmod +x /backend/ddalab-manager

# Expose port 8080 for the backend
EXPOSE 8080

WORKDIR /

# Docker Extension specific labels - these are REQUIRED for Docker Hub to recognize as extension
LABEL org.opencontainers.image.title="DDALAB Manager" \
      org.opencontainers.image.description="Manage DDALAB (Delay Differential Analysis Laboratory) installations" \
      org.opencontainers.image.vendor="DDALAB Team" \
      org.opencontainers.image.source="https://github.com/sdraeger/DDALAB" \
      org.opencontainers.image.url="https://github.com/sdraeger/DDALAB" \
      org.opencontainers.image.documentation="https://github.com/sdraeger/DDALAB/tree/main/docker-extension" \
      org.opencontainers.image.licenses="MIT" \
      com.docker.desktop.extension.api.version="0.3.4" \
      com.docker.desktop.extension.icon="https://avatars.githubusercontent.com/u/5429470?s=200&v=4" \
      com.docker.extension.categories="utility,data-science" \
      com.docker.extension.screenshots="[]" \
      com.docker.extension.detailed-description="DDALAB Manager is a Docker Desktop extension that helps you manage DDALAB (Delay Differential Analysis Laboratory) installations. Monitor services, view logs, create backups, and control your DDALAB deployment directly from Docker Desktop's Extensions tab." \
      com.docker.extension.publisher-url="https://github.com/sdraeger/DDALAB" \
      com.docker.extension.additional-urls='[{"title":"Documentation","url":"https://github.com/sdraeger/DDALAB/blob/main/README.md"},{"title":"Issues","url":"https://github.com/sdraeger/DDALAB/issues"}]' \
      com.docker.extension.changelog="Initial release of DDALAB Manager extension with service management, monitoring, and configuration capabilities."

CMD ["/backend/ddalab-manager"]