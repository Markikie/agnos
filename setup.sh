#!/bin/bash

echo "Setting up Hospital A API with HTTPS..."

# Create necessary directories
mkdir -p ssl logs/nginx

# Generate SSL certificate
echo "Generating self-signed SSL certificate..."
bash scripts/generate-ssl.sh

# Build and start services
echo "Building and starting Docker services..."
docker-compose down
docker-compose build
docker-compose up -d

echo ""
echo "Setup complete!"
echo ""
echo "Services:"
echo "- Go App: http://localhost:8080"
echo "- Nginx HTTP: http://localhost:80"
echo "- Nginx HTTPS: https://localhost:443"
echo ""
echo "To test with the domain name, add this to your hosts file:"
echo "127.0.0.1 hospital-a.api.co.th"
echo ""
echo "Then visit: https://hospital-a.api.co.th"
echo ""
echo "Note: You'll need to accept the self-signed certificate in your browser."
