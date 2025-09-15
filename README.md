# DDALAB Docker Desktop Extension

A Docker Desktop extension for managing DDALAB (Delay Differential Analysis Laboratory) installations.

## Features

- 🚀 **Service Management**: Start, stop, and restart individual DDALAB services
- 📊 **Status Monitoring**: Real-time status of all DDALAB components
- 📝 **Log Viewer**: View recent logs from all services
- 💾 **Backup Management**: Create database backups with one click
- 🎯 **Stack Control**: Manage the entire DDALAB stack easily

## Prerequisites

- Docker Desktop 4.8.0 or later
- An existing DDALAB installation (from DDALAB-setup)
- Docker Compose installed on the host

## Installation

### From Docker Hub

```bash
docker extension install sdraeger/ddalab-manager:latest
```

### Build from Source

```bash
# Clone this repository
git clone https://github.com/sdraeger/DDALAB.git
cd DDALAB/docker-extension

# Build the extension
make build

# Install the extension
make install
```

## Usage

1. Open Docker Desktop
2. Navigate to the Extensions tab
3. Click on "DDALAB Manager"
4. The extension will automatically detect your DDALAB installation

### Service Management

- Click on individual service controls to start, stop, or restart services
- Use the "Stack Actions" to control all services at once

### Viewing Logs

- Click "Refresh Logs" to load the latest logs
- Logs are displayed in a scrollable viewer
- Shows the last 100 lines from all services

### Creating Backups

- Click "Create Backup" to create a PostgreSQL database backup
- Backups are stored in the `backups` directory of your DDALAB installation

## Development

### Project Structure

```
docker-extension/
├── backend/          # Go backend for Docker API interactions
├── ui/              # Frontend React application
├── vm/              # Docker Compose configuration for the extension
├── metadata.json    # Extension metadata
├── extension.yaml   # Extension configuration
├── Dockerfile       # Multi-stage build for the extension
└── Makefile        # Build and development commands
```

### Development Commands

```bash
# Build the extension
make build

# Install for development
make install

# Update after changes
make update

# Enable debug mode
make dev

# View logs
make logs

# Validate extension
make validate

# Uninstall
make uninstall
```

### Backend API Endpoints

- `GET /api/status` - Get current status of DDALAB services
- `POST /api/services/{service}/{action}` - Control individual services
- `POST /api/stack/{action}` - Control the entire stack
- `GET /api/logs` - Retrieve recent logs
- `POST /api/backup` - Create a database backup

## Troubleshooting

### Extension doesn't detect DDALAB installation

The extension looks for DDALAB-setup in these locations:
- `/DDALAB-setup`
- `../DDALAB-setup`
- `~/DDALAB-setup`
- `~/Desktop/DDALAB-setup`

Make sure your installation is in one of these paths or create a symlink.

### Services show as "stopped" when they're running

Ensure the Docker socket is properly mounted and accessible. The extension needs access to Docker to query container status.

### Backup fails

- Verify PostgreSQL container is running
- Check that the `ddalab` database exists
- Ensure sufficient disk space for backups

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This extension is part of the DDALAB project and follows the same license terms.

## Support

For issues and questions:
- GitHub Issues: [github.com/sdraeger/DDALAB/issues](https://github.com/sdraeger/DDALAB/issues)
- Email: sdraeger@salk.edu