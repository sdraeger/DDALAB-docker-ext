#!/bin/bash

# DDALAB Docker Extension Backend Startup Script

echo "Starting DDALAB Docker Extension Backend..."

# Stop any existing backend
docker stop ddalab-extension-backend 2>/dev/null
docker rm ddalab-extension-backend 2>/dev/null

# Start the backend with proper mounts
docker run -d \
  --name ddalab-extension-backend \
  --restart unless-stopped \
  -p 8080:8080 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /Users/simon/Desktop/DDALAB-setup:/Users/simon/Desktop/DDALAB-setup:ro \
  sdraeger1/ddalab-docker-ext:latest

# Wait a moment for startup
sleep 3

# Check if backend is running
if docker ps | grep -q ddalab-extension-backend; then
    echo "✅ Backend started successfully!"
    echo "📍 Backend URL: http://localhost:8080"
    echo "🔍 Testing API..."
    
    # Test the API
    if curl -s http://localhost:8080/api/status > /dev/null; then
        echo "✅ API is responding!"
        echo "🎯 You can now use the Docker Desktop extension"
    else
        echo "❌ API is not responding yet, please wait a moment"
    fi
else
    echo "❌ Failed to start backend"
    docker logs ddalab-extension-backend
fi

echo ""
echo "To stop the backend:"
echo "  docker stop ddalab-extension-backend"
echo ""
echo "To view logs:"
echo "  docker logs -f ddalab-extension-backend"