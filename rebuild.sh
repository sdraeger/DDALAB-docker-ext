#!/bin/bash
# Script to rebuild and update the DDALAB Docker extension

set -e  # Exit on error

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
EXTENSION_NAME="sdraeger1/ddalab-docker-ext:latest"

echo "🔨 Building DDALAB Docker Extension..."
echo "======================================"

# Change to extension directory
cd "$SCRIPT_DIR"

# Build the UI
echo "📦 Building UI..."
cd ui
npm run build
if [ $? -ne 0 ]; then
    echo "❌ UI build failed!"
    exit 1
fi

# Copy built files to root
cp dist/bundle.js ..
cp dist/index.html ..
cd ..

echo "✅ UI build complete"

# Build the Docker image
echo "🐳 Building Docker image..."
docker build -t "$EXTENSION_NAME" .
if [ $? -ne 0 ]; then
    echo "❌ Docker build failed!"
    exit 1
fi

echo "✅ Docker image built successfully"

# Check if extension is already installed
if docker extension ls | grep -q "sdraeger1/ddalab-docker-ext"; then
    echo "🔄 Updating existing extension..."
    echo "y" | docker extension update "$EXTENSION_NAME"
else
    echo "📥 Installing extension..."
    echo "y" | docker extension install "$EXTENSION_NAME"
fi

if [ $? -eq 0 ]; then
    echo ""
    echo "✨ Extension successfully rebuilt and updated!"
    echo ""
    echo "You can now use the DDALAB Manager in Docker Desktop."
    echo "Check the extension tab in Docker Desktop to access it."
else
    echo ""
    echo "❌ Extension update/install failed!"
    echo "Please check the error messages above."
    exit 1
fi