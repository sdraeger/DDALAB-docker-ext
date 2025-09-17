#!/bin/bash
# Development script to rebuild and update the DDALAB Docker extension with options

set -e  # Exit on error

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
EXTENSION_NAME="sdraeger/ddalab-manager:latest"

# Default options
SKIP_UI=false
SKIP_BACKEND=false
VERBOSE=false
FORCE_REINSTALL=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --skip-ui)
            SKIP_UI=true
            shift
            ;;
        --skip-backend)
            SKIP_BACKEND=true
            shift
            ;;
        --verbose|-v)
            VERBOSE=true
            shift
            ;;
        --force|-f)
            FORCE_REINSTALL=true
            shift
            ;;
        --help|-h)
            echo "Usage: $0 [options]"
            echo ""
            echo "Options:"
            echo "  --skip-ui       Skip UI rebuild (use existing bundle.js)"
            echo "  --skip-backend  Skip backend rebuild"
            echo "  --verbose, -v   Show detailed output"
            echo "  --force, -f     Force reinstall (remove and install fresh)"
            echo "  --help, -h      Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0                    # Full rebuild"
            echo "  $0 --skip-ui          # Rebuild only backend"
            echo "  $0 --skip-backend     # Rebuild only UI"
            echo "  $0 --force            # Remove and reinstall extension"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

echo "🔨 Building DDALAB Docker Extension"
echo "==================================="
echo ""

# Change to extension directory
cd "$SCRIPT_DIR"

# Build the UI if not skipped
if [ "$SKIP_UI" = false ]; then
    echo "📦 Building UI..."
    cd ui
    
    if [ "$VERBOSE" = true ]; then
        npm run build
    else
        npm run build > /dev/null 2>&1
    fi
    
    if [ $? -ne 0 ]; then
        echo "❌ UI build failed!"
        exit 1
    fi
    
    # Copy built files to root
    cp dist/bundle.js ..
    cp dist/index.html ..
    cd ..
    
    echo "✅ UI build complete"
else
    echo "⏭️  Skipping UI build"
fi

# Build the Docker image
echo "🐳 Building Docker image..."

BUILD_ARGS=""
if [ "$SKIP_BACKEND" = true ]; then
    BUILD_ARGS="--target stage-1"
    echo "⏭️  Skipping backend compilation"
fi

if [ "$VERBOSE" = true ]; then
    docker build $BUILD_ARGS -t "$EXTENSION_NAME" .
else
    docker build $BUILD_ARGS -t "$EXTENSION_NAME" . > /dev/null
fi

if [ $? -ne 0 ]; then
    echo "❌ Docker build failed!"
    exit 1
fi

echo "✅ Docker image built successfully"

# Force reinstall if requested
if [ "$FORCE_REINSTALL" = true ]; then
    if docker extension ls | grep -q "sdraeger/ddalab-manager"; then
        echo "🗑️  Removing existing extension..."
        docker extension rm sdraeger/ddalab-manager 2>/dev/null || true
    fi
fi

# Check if extension is already installed
if docker extension ls | grep -q "sdraeger/ddalab-manager"; then
    echo "🔄 Updating existing extension..."
    if [ "$VERBOSE" = true ]; then
        echo "y" | docker extension update "$EXTENSION_NAME"
    else
        echo "y" | docker extension update "$EXTENSION_NAME" > /dev/null
    fi
else
    echo "📥 Installing extension..."
    if [ "$VERBOSE" = true ]; then
        echo "y" | docker extension install "$EXTENSION_NAME"
    else
        echo "y" | docker extension install "$EXTENSION_NAME" > /dev/null
    fi
fi

if [ $? -eq 0 ]; then
    echo ""
    echo "✨ Extension successfully rebuilt and updated!"
    echo ""
    
    # Show extension status
    echo "📊 Extension Status:"
    docker extension ls | grep -E "(PROVIDER|sdraeger/ddalab-manager)" | column -t
    
    echo ""
    echo "🚀 The DDALAB Manager is ready in Docker Desktop!"
else
    echo ""
    echo "❌ Extension update/install failed!"
    echo "Try running with --verbose flag for more details"
    exit 1
fi