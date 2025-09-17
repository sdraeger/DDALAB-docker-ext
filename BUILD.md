# Building and Updating the DDALAB Docker Extension

This document describes how to build and update the DDALAB Docker Desktop extension.

**Note**: The UI is now written in TypeScript for better type safety and development experience.

## Prerequisites

- Docker Desktop installed and running
- Node.js and npm installed
- Go 1.21+ installed (for backend development)

## Quick Start

### Using Make (Recommended)

```bash
# Full rebuild and update
make

# Or explicitly
make rebuild
```

### Using Scripts

```bash
# Full rebuild and update
./rebuild.sh

# Development rebuild with options
./rebuild-dev.sh --help
```

## Available Commands

### Make Targets

| Command | Description |
|---------|-------------|
| `make` | Full rebuild and update (default) |
| `make build` | Build UI and Docker image |
| `make install` | Build and install extension (first time) |
| `make update` | Build and update existing extension |
| `make rebuild` | Full rebuild and update |
| `make backend` | Rebuild only backend |
| `make frontend` | Rebuild only frontend |
| `make remove` | Remove the extension |
| `make clean` | Clean build artifacts |
| `make help` | Show help message |
| `npm run typecheck` | Run TypeScript type checking (in ui directory) |
| `npm run lint` | Run ESLint on TypeScript files (in ui directory) |

### Development Options

```bash
# Verbose output
make dev-rebuild ARGS='--verbose'

# Force reinstall (remove and install fresh)
make dev-rebuild ARGS='--force'

# Skip UI build (backend only)
make dev-rebuild ARGS='--skip-ui'

# Skip backend build (UI only)
make dev-rebuild ARGS='--skip-backend'
```

### Script Options

The `rebuild-dev.sh` script provides fine-grained control:

```bash
# Show all options
./rebuild-dev.sh --help

# Examples
./rebuild-dev.sh --skip-ui          # Rebuild only backend
./rebuild-dev.sh --skip-backend     # Rebuild only UI
./rebuild-dev.sh --verbose          # Show detailed output
./rebuild-dev.sh --force            # Remove and reinstall
```

## Development Workflow

### Frontend Changes Only

```bash
make frontend
```

This is the fastest way to test UI changes without rebuilding the backend.

### Backend Changes Only

```bash
make backend
```

This skips the UI build and only rebuilds the Go backend.

### Full Rebuild

```bash
make
```

Use this when you've made changes to both frontend and backend, or when updating dependencies.

## Troubleshooting

### Extension Not Updating

If the extension doesn't seem to update:

1. Force reinstall:
   ```bash
   make dev-rebuild ARGS='--force'
   ```

2. Or manually remove and reinstall:
   ```bash
   make remove
   make install
   ```

### Build Failures

For detailed error messages:

```bash
make dev-rebuild ARGS='--verbose'
```

### Checking Extension Status

```bash
docker extension ls | grep ddalab
```

## File Structure

- `rebuild.sh` - Simple rebuild script
- `rebuild-dev.sh` - Advanced rebuild script with options
- `Makefile` - Make targets for common tasks
- `ui/` - Frontend React application
- `backend/` - Go backend service
- `Dockerfile` - Multi-stage build for the extension
- `compose.yaml` - Extension service configuration
- `metadata.json` - Extension metadata