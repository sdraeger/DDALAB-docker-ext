# Publishing DDALAB Docker Extension

This document explains how to properly publish the DDALAB Docker Extension to Docker Hub so it appears as an extension rather than a regular container image.

## Key Requirements for Docker Extensions

For Docker Hub to recognize an image as a Docker Desktop extension, it must:

1. **Have the correct OCI labels** (✅ now included in Dockerfile)
2. **Include extension metadata files** (✅ now included)
3. **Be tagged and pushed correctly** 
4. **Have the right repository structure**

## Current Extension Structure

```
docker-extension/
├── Dockerfile                    # Multi-stage build with extension labels
├── metadata.json                 # Extension metadata for Docker Desktop
├── extension.yaml               # Extension configuration
├── extension-manifest.json     # Complete extension manifest
├── compose.yaml                 # Service composition for extension
├── icon.svg                     # Extension icon
├── backend/                     # Go backend API
└── ui/                         # React frontend UI
```

## Required Labels (Now Included)

The Dockerfile now includes these essential labels:

```dockerfile
LABEL org.opencontainers.image.title="DDALAB Manager" \
      org.opencontainers.image.description="Manage DDALAB installations" \
      org.opencontainers.image.vendor="DDALAB Team" \
      org.opencontainers.image.source="https://github.com/sdraeger/DDALAB" \
      com.docker.desktop.extension.api.version="0.3.4" \
      com.docker.extension.categories="utility,data-science,scientific-computing"
```

## Publishing Steps

### 1. Build the Extension

```bash
cd docker-extension
docker build -t sdraeger1/ddalab-docker-ext:latest .
```

### 2. Tag for Docker Hub

```bash
# Tag as extension (important!)
docker tag sdraeger1/ddalab-docker-ext:latest sdraeger1/ddalab-docker-ext:1.0.0
docker tag sdraeger1/ddalab-docker-ext:latest sdraeger1/ddalab-docker-ext:latest
```

### 3. Push to Docker Hub

```bash
docker push sdraeger1/ddalab-docker-ext:1.0.0
docker push sdraeger1/ddalab-docker-ext:latest
```

### 4. Verify Extension Labels

```bash
# Check that the image has the correct labels
docker inspect sdraeger1/ddalab-docker-ext:latest | jq '.[0].Config.Labels'
```

You should see labels like:
- `com.docker.desktop.extension.api.version`
- `com.docker.extension.categories`
- `org.opencontainers.image.title`

## Docker Hub Repository Settings

On Docker Hub, ensure:

1. **Repository Description**: "Docker Desktop extension for managing DDALAB installations"
2. **README**: Should explain it's a Docker Desktop extension
3. **Tags**: Use semantic versioning (1.0.0, 1.0.1, etc.)

## Testing the Extension

### Local Testing

```bash
# Install extension locally
docker extension install sdraeger1/ddalab-docker-ext:latest

# Or install from local build
docker extension install .

# Check extension status
docker extension ls

# Remove for testing
docker extension remove sdraeger1/ddalab-docker-ext
```

### Verify in Docker Desktop

1. Open Docker Desktop
2. Go to Extensions tab
3. Look for "DDALAB Manager" in the installed extensions
4. It should appear with the proper icon and description

## Troubleshooting

### If Extension Appears as Regular Image

This means Docker Hub doesn't recognize it as an extension. Check:

1. **Labels are present**: `docker inspect` should show extension labels
2. **Metadata files exist**: metadata.json should be in the image
3. **API version is current**: Should be "0.3.4" or later
4. **Categories are valid**: Use standard Docker extension categories

### Common Issues

1. **Missing `com.docker.desktop.extension.api.version` label**
2. **Incorrect metadata.json format**
3. **Missing or malformed icon**
4. **Wrong repository tagging**

## Extension Categories

Valid categories for Docker extensions:
- `utility` - General utility tools
- `data-science` - Data science and analytics
- `scientific-computing` - Scientific computing tools
- `development` - Development tools
- `monitoring` - Monitoring and observability
- `security` - Security tools

## Next Steps

After publishing with the updated configuration:

1. **Re-push the image** with the new labels and metadata
2. **Wait 10-15 minutes** for Docker Hub to process the changes
3. **Check the Docker Hub repository** - it should show as an extension
4. **Test installation** in Docker Desktop
5. **Update launcher bootstrap** to use the correct extension name

The extension should now appear properly as a Docker Desktop extension rather than a regular container image.