# Makefile for DDALAB Docker Extension

.PHONY: build build-ui build-docker install update remove clean rebuild dev-rebuild help

# Variables
EXTENSION_NAME = sdraeger1/ddalab-manager:latest
UI_DIR = ui
DIST_FILES = bundle.js index.html

# Default target
all: rebuild

# Build UI only
build-ui:
	@echo "📦 Building UI..."
	@cd $(UI_DIR) && npm run build
	@cp $(UI_DIR)/dist/bundle.js .
	@cp $(UI_DIR)/dist/index.html .
	@echo "✅ UI build complete"

# Build Docker image only
build-docker:
	@echo "🐳 Building Docker image..."
	@docker build -t $(EXTENSION_NAME) .
	@echo "✅ Docker image built"

# Build everything
build: build-ui build-docker

# Install extension (first time)
install: build
	@echo "📥 Installing extension..."
	@echo "y" | docker extension install $(EXTENSION_NAME)

# Update existing extension
update: build
	@echo "🔄 Updating extension..."
	@echo "y" | docker extension update $(EXTENSION_NAME)

# Remove extension
remove:
	@echo "🗑️  Removing extension..."
	@docker extension rm sdraeger1/ddalab-manager || true
	@docker extension rm sdraeger/ddalab-manager || true

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -f bundle.js index.html
	@rm -rf $(UI_DIR)/dist
	@echo "✅ Clean complete"

# Full rebuild and update
rebuild:
	@./rebuild.sh

# Development rebuild with options
dev-rebuild:
	@./rebuild-dev.sh $(ARGS)

# Quick backend-only rebuild
backend: 
	@./rebuild-dev.sh --skip-ui

# Quick frontend-only rebuild
frontend:
	@./rebuild-dev.sh --skip-backend

# Show help
help:
	@echo "DDALAB Docker Extension - Make targets"
	@echo ""
	@echo "Available targets:"
	@echo "  make build          - Build UI and Docker image"
	@echo "  make install        - Build and install extension (first time)"
	@echo "  make update         - Build and update existing extension"
	@echo "  make rebuild        - Full rebuild and update (default)"
	@echo "  make backend        - Rebuild only backend"
	@echo "  make frontend       - Rebuild only frontend"
	@echo "  make remove         - Remove the extension"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make help           - Show this help message"
	@echo ""
	@echo "Development options:"
	@echo "  make dev-rebuild ARGS='--verbose'     - Verbose rebuild"
	@echo "  make dev-rebuild ARGS='--force'       - Force reinstall"
	@echo "  make dev-rebuild ARGS='--skip-ui'     - Skip UI build"
	@echo ""