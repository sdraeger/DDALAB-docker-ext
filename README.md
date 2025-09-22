# DDALAB Manager - Docker Desktop Extension

[![Docker Hub](https://img.shields.io/docker/v/sdraeger1/ddalab-docker-ext?label=Docker%20Hub)](https://hub.docker.com/r/sdraeger1/ddalab-docker-ext)
[![Docker Image Size](https://img.shields.io/docker/image-size/sdraeger1/ddalab-docker-ext/latest)](https://hub.docker.com/r/sdraeger1/ddalab-docker-ext)
[![Docker Pulls](https://img.shields.io/docker/pulls/sdraeger1/ddalab-docker-ext)](https://hub.docker.com/r/sdraeger1/ddalab-docker-ext)

A Docker Desktop extension for managing DDALAB (Delay Differential Analysis Laboratory) installations directly from Docker Desktop's Extensions tab.

## 🧪 What is DDALAB?

DDALAB (Delay Differential Analysis Laboratory) is a comprehensive scientific computing platform for performing Delay Differential Analysis on EDF (European Data Format) and ASCII files. It is planned to be used in neuroscience research for analyzing brain activity patterns and detecting anomalies in EEG and other physiological data.

## 🚀 Features

### 📊 **Service Management**

- **Real-time Status Monitoring** - View the health and status of all DDALAB services
- **One-click Service Control** - Start, stop, and restart DDALAB stack with ease
- **Resource Monitoring** - Track CPU, memory, and uptime for each service

### 📋 **Log Management**

- **Live Log Streaming** - View real-time logs from all DDALAB services
- **Log Filtering** - Filter logs by service, log level, or time range
- **Log Export** - Download logs for debugging and analysis

### 🔧 **Configuration Management**

- **Environment Configuration** - Manage DDALAB environment variables
- **Service Configuration** - Configure individual service settings
- **Backup & Restore** - Create and restore configuration backups

### 🐳 **Docker Integration**

- **Native Docker Desktop Integration** - Seamlessly integrated into Docker Desktop UI
- **Container Management** - Direct control over DDALAB containers
- **Network Monitoring** - View container networking and port mappings

## 📦 Installation

### Prerequisites

- **Docker Desktop 4.10.0+** with Extensions support
- **Docker Engine 20.10+**
- **4GB RAM minimum** (8GB recommended for DDALAB)
- **10GB free disk space** for DDALAB data and containers

### Install from Docker Desktop

1. Open **Docker Desktop**
2. Navigate to the **Extensions** tab
3. Search for **"DDALAB Manager"**
4. Click **Install**

### Install from Command Line

```bash
# Install the extension
docker extension install sdraeger1/ddalab-docker-ext:latest

# Verify installation
docker extension ls
```

### Manual Installation

```bash
# Clone the repository
git clone https://github.com/sdraeger/DDALAB.git
cd DDALAB/docker-extension

# Build and install locally
docker extension install .
```

## 🎯 Quick Start

### 1. **Launch DDALAB Manager**

After installation, find "DDALAB Manager" in your Docker Desktop Extensions tab.

### 2. **Initialize DDALAB**

- Click **"Initialize DDALAB"** to set up the required services
- The extension will automatically pull and configure:
  - DDALAB API Server (Python FastAPI)
  - PostgreSQL Database
  - Redis Cache
  - MinIO Object Storage
  - Traefik Reverse Proxy

### 3. **Monitor Services**

- View real-time status of all DDALAB components
- Monitor resource usage and health checks
- Access service logs and metrics

### 4. **Access DDALAB**

Once services are running, access DDALAB at:

- **Main Application**: https://localhost (via Traefik)
- **API Documentation**: https://localhost/docs
- **MinIO Console**: https://localhost/minio

## 🏗️ Architecture

DDALAB Manager orchestrates a complete microservices stack:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   React UI      │    │  DDALAB API     │    │  PostgreSQL     │
│   (Frontend)    │◄──►│  (FastAPI)      │◄──►│  (Database)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
        │                        │                        │
        ▼                        ▼                        ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Traefik       │    │     Redis       │    │     MinIO       │
│   (Proxy)       │    │    (Cache)      │    │   (Storage)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Services Included

| Service      | Purpose                      | Port      | Health Check |
| ------------ | ---------------------------- | --------- | ------------ |
| **ddalab**   | Main API server with GraphQL | 8000      | ✅ Built-in  |
| **postgres** | Primary database             | 5432      | ✅ Built-in  |
| **redis**    | Caching and sessions         | 6379      | ✅ Built-in  |
| **minio**    | Object storage for EDF files | 9000/9001 | ✅ Built-in  |
| **traefik**  | Reverse proxy with SSL       | 80/443    | ✅ Built-in  |

## 📋 Usage Examples

### Basic Service Management

```bash
# Start DDALAB stack
curl -X POST http://localhost:8080/api/v1/start

# Check service status
curl http://localhost:8080/api/v1/status

# View logs
curl http://localhost:8080/api/v1/logs/ddalab

# Stop services
curl -X POST http://localhost:8080/api/v1/stop
```

### Configuration Management

The extension provides a user-friendly interface for:

- Managing environment variables
- Configuring service settings
- Creating configuration backups
- Monitoring resource usage

### Advanced Features

- **Custom Deployment Profiles** - Switch between development and production configurations
- **Automated Backups** - Schedule regular backups of DDALAB data
- **Performance Monitoring** - Track system performance and resource usage
- **Log Analysis** - Built-in log parsing and analysis tools

## 🔧 Configuration

### Environment Variables

The extension supports configuration through environment variables:

```env
DDALAB_MODE=production          # deployment mode
DDALAB_DEBUG=false             # debug logging
DDALAB_HOST=localhost          # host domain
DDALAB_PORT=443                # HTTPS port
POSTGRES_PASSWORD=secure123    # database password
REDIS_PASSWORD=secure456       # Redis password
MINIO_ROOT_USER=admin         # MinIO admin user
MINIO_ROOT_PASSWORD=secure789  # MinIO admin password
```

### Volume Mounts

DDALAB requires persistent storage for:

- **Database**: PostgreSQL data
- **Object Storage**: EDF files and analysis results
- **Configuration**: Service configurations and SSL certificates
- **Logs**: Application and service logs

### Network Configuration

The extension creates a dedicated Docker network for DDALAB services with:

- **Internal communication** between services
- **External access** via Traefik proxy
- **SSL termination** for secure connections
- **Load balancing** for high availability

## 🛠️ Development

### Building from Source

```bash
# Clone the repository
git clone https://github.com/sdraeger/DDALAB.git
cd DDALAB/docker-extension

# Build the UI
cd ui
npm install
npm run build
cd ..

# Build the extension
docker build -t ddalab-docker-ext .

# Install locally
docker extension install ddalab-docker-ext:latest
```

### Extension Structure

```
docker-extension/
├── Dockerfile              # Multi-stage build configuration
├── compose.yaml            # Service composition
├── metadata.json           # Extension metadata
├── extension.yaml          # Extension configuration
├── icon.svg               # Extension icon
├── backend/               # Go backend API
│   ├── main.go
│   ├── handlers/
│   └── models/
└── ui/                    # React frontend
    ├── src/
    ├── package.json
    └── webpack.config.js
```

### API Endpoints

The extension backend provides RESTful API endpoints:

- `GET /api/v1/health` - Health check
- `GET /api/v1/status` - Service status
- `POST /api/v1/start` - Start services
- `POST /api/v1/stop` - Stop services
- `GET /api/v1/logs/{service}` - Get service logs
- `GET /api/v1/config` - Get configuration
- `POST /api/v1/config` - Update configuration

## 🧪 Testing

### Automated Testing

```bash
# Run backend tests
cd backend
go test ./...

# Run frontend tests
cd ui
npm test

# Run integration tests
npm run test:integration
```

### Manual Testing

1. **Install the extension** in Docker Desktop
2. **Launch DDALAB Manager** from Extensions tab
3. **Initialize services** and verify all components start
4. **Test service management** - start, stop, restart operations
5. **Verify logging** - check that logs are accessible and formatted correctly
6. **Test configuration** - modify settings and verify they persist

## 📚 Documentation

### Additional Resources

- **[DDALAB Documentation](https://github.com/sdraeger/DDALAB/blob/main/README.md)** - Complete DDALAB user guide
- **[API Reference](https://github.com/sdraeger/DDALAB/tree/main/packages/api)** - DDALAB API documentation
- **[Development Guide](https://github.com/sdraeger/DDALAB/blob/main/docs/DEVELOPMENT_SETUP.md)** - Setup development environment
- **[Deployment Guide](https://github.com/sdraeger/DDALAB/blob/main/docs/DEPLOYMENT.md)** - Production deployment instructions

### Scientific Background

DDALAB implements advanced algorithms for:

- **Delay Differential Analysis (DDA)** - Mathematical framework for analyzing time-delayed systems
- **EDF File Processing** - European Data Format for biological signal storage
- **Real-time Signal Processing** - Live analysis of neurological data
- **Machine Learning Integration** - AI-powered pattern recognition

## 🐛 Troubleshooting

### Common Issues

#### Extension Won't Start

```bash
# Check Docker Desktop version
docker version

# Verify extension installation
docker extension ls

# Check extension logs
docker extension logs sdraeger1/ddalab-docker-ext
```

#### Services Fail to Start

```bash
# Check available resources
docker system df

# Verify port availability
netstat -an | grep :8080

# Check service logs
curl http://localhost:8080/api/v1/logs/ddalab
```

#### Performance Issues

- **Increase Docker Desktop memory allocation** (8GB+ recommended)
- **Close other resource-intensive applications**
- **Check disk space** (DDALAB requires 10GB+ free space)
- **Restart Docker Desktop** if services become unresponsive

### Support Channels

- **[GitHub Issues](https://github.com/sdraeger/DDALAB/issues)** - Bug reports and feature requests
- **[Discussions](https://github.com/sdraeger/DDALAB/discussions)** - Community support and questions
- **[Documentation](https://github.com/sdraeger/DDALAB/tree/main/docs)** - Comprehensive guides and tutorials

## 🤝 Contributing

We welcome contributions to DDALAB Manager! Please see our [Contributing Guide](https://github.com/sdraeger/DDALAB/blob/main/CONTRIBUTING.md) for details.

### Development Setup

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Make your changes** and add tests
4. **Test thoroughly** using the testing procedures above
5. **Submit a pull request** with a clear description

### Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/version/2/1/code_of_conduct/).

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/sdraeger/DDALAB/blob/main/LICENSE) file for details.

## 🙏 Acknowledgments

- **Docker Desktop Extensions Team** - For providing the extension framework
- **DDALAB Community** - For feedback and contributions
- **Scientific Computing Community** - For advancing delay differential analysis research

---

**Made with ❤️ for the scientific computing community**

For questions, suggestions, or support, please [open an issue](https://github.com/sdraeger/DDALAB/issues) on GitHub.
