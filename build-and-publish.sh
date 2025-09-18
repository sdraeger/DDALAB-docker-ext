#!/bin/bash

# Build and publish DDALAB Docker Extension
# This script ensures the extension is built and tagged correctly for Docker Hub

set -e

# Configuration
REGISTRY="sdraeger1"
IMAGE_NAME="ddalab-docker-ext" 
VERSION="1.0.0"
FULL_IMAGE="$REGISTRY/$IMAGE_NAME"

echo "🔨 Building DDALAB Docker Extension..."
echo "Registry: $REGISTRY"
echo "Image: $IMAGE_NAME"
echo "Version: $VERSION"
echo "Full image name: $FULL_IMAGE"

# Build the UI first
echo "📦 Building UI..."
cd ui
npm install
npm run build
cd ..

# Build the extension image
echo "🐳 Building Docker extension image..."
docker build -t "$FULL_IMAGE:$VERSION" -t "$FULL_IMAGE:latest" .

# Verify the image has correct extension labels
echo "🔍 Verifying extension labels..."
if docker inspect "$FULL_IMAGE:latest" | jq '.[0].Config.Labels' | grep -q "com.docker.desktop.extension.api.version"; then
    echo "✅ Extension labels found!"
else
    echo "❌ Extension labels missing! The image may not be recognized as an extension."
    exit 1
fi

# Show the labels for verification
echo "📋 Extension labels:"
docker inspect "$FULL_IMAGE:latest" | jq '.[0].Config.Labels | with_entries(select(.key | startswith("com.docker.desktop.extension") or startswith("com.docker.extension") or startswith("org.opencontainers.image")))'

# Test locally first
echo "🧪 Testing extension locally..."
echo "Removing any existing extension..."
docker extension remove "$FULL_IMAGE" 2>/dev/null || true

echo "Installing extension locally for testing..."
if docker extension install "$FULL_IMAGE:latest"; then
    echo "✅ Extension installed successfully!"
    echo "Check Docker Desktop Extensions tab to verify it appears correctly."
    echo ""
    echo "Extension info:"
    docker extension ls | grep "$IMAGE_NAME" || echo "Extension not found in list"
else
    echo "❌ Extension installation failed!"
    exit 1
fi

# Ask for confirmation before pushing
echo ""
read -p "🚀 Push to Docker Hub? (y/N): " -n 1 -r
echo ""

if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "📤 Pushing to Docker Hub..."
    
    # Login check
    if ! docker info | grep -q Username; then
        echo "Please login to Docker Hub first:"
        docker login
    fi
    
    # Push both tags
    docker push "$FULL_IMAGE:$VERSION"
    docker push "$FULL_IMAGE:latest"
    
    echo "✅ Extension published successfully!"
    echo ""
    echo "📝 Next steps:"
    echo "1. Check Docker Hub: https://hub.docker.com/r/$REGISTRY/$IMAGE_NAME"
    echo "2. Verify it appears as an extension (not just a container image)"
    echo "3. Wait 10-15 minutes for Docker Hub to process the metadata"
    echo "4. Test installation from Docker Hub: docker extension install $FULL_IMAGE:latest"
    echo ""
    echo "🔍 Extension details:"
    echo "- Name: DDALAB Manager"
    echo "- Image: $FULL_IMAGE:$VERSION"
    echo "- Categories: utility, data-science, scientific-computing"
    echo "- API Version: 0.3.4"
else
    echo "Push cancelled. Extension is built and tested locally."
fi

echo ""
echo "🎉 Build process complete!"