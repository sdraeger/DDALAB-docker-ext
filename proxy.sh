#!/bin/sh

# Start a simple HTTP server that serves UI and proxies API calls
cat > /tmp/proxy.py << 'EOF'
#!/usr/bin/env python3
import http.server
import socketserver
import urllib.request
import urllib.error
import json
import os
from urllib.parse import urlparse

class ProxyHandler(http.server.SimpleHTTPRequestHandler):
    def __init__(self, *args, **kwargs):
        # Set the directory to serve from
        super().__init__(*args, directory="/ui", **kwargs)
    
    def do_GET(self):
        if self.path.startswith('/api/'):
            self.proxy_request('GET')
        else:
            # Serve static files from /ui
            super().do_GET()
    
    def do_POST(self):
        if self.path.startswith('/api/'):
            self.proxy_request('POST')
        else:
            self.send_error(404)
    
    def do_PUT(self):
        if self.path.startswith('/api/'):
            self.proxy_request('PUT')
        else:
            self.send_error(404)
    
    def proxy_request(self, method):
        # Forward request to ddalab-control
        backend_url = f"http://ddalab-control:8080{self.path}"
        
        try:
            # Read request body if present
            content_length = int(self.headers.get('Content-Length', 0))
            request_body = self.rfile.read(content_length) if content_length > 0 else None
            
            # Create request
            req = urllib.request.Request(backend_url, data=request_body, method=method)
            
            # Copy headers (excluding some that shouldn't be forwarded)
            for header_name, header_value in self.headers.items():
                if header_name.lower() not in ['host', 'content-length']:
                    req.add_header(header_name, header_value)
            
            # Make request to backend
            with urllib.request.urlopen(req, timeout=30) as response:
                # Send response headers
                self.send_response(response.status)
                for header_name, header_value in response.headers.items():
                    if header_name.lower() not in ['content-length', 'transfer-encoding']:
                        self.send_header(header_name, header_value)
                
                # Read and send response body
                response_body = response.read()
                self.send_header('Content-Length', str(len(response_body)))
                self.end_headers()
                self.wfile.write(response_body)
                
        except urllib.error.HTTPError as e:
            self.send_response(e.code)
            self.send_header('Content-Type', 'application/json')
            self.end_headers()
            error_response = json.dumps({"error": f"Backend error: {e.reason}"}).encode()
            self.wfile.write(error_response)
            
        except Exception as e:
            self.send_response(500)
            self.send_header('Content-Type', 'application/json')
            self.end_headers()
            error_response = json.dumps({"error": f"Proxy error: {str(e)}"}).encode()
            self.wfile.write(error_response)

# Start server
PORT = 8080
with socketserver.TCPServer(("", PORT), ProxyHandler) as httpd:
    print(f"Serving UI at port {PORT} with API proxy to ddalab-control:8080")
    httpd.serve_forever()
EOF

# Install python and start the proxy
apk add --no-cache python3
python3 /tmp/proxy.py