# Agnos - Hospital A API

A Go-based API service with nginx reverse proxy and HTTPS support.

## Features

- Go web application with Gin framework
- Nginx reverse proxy with SSL termination
- Self-signed SSL certificate for HTTPS
- PostgreSQL database
- Docker containerization

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Bash (for running setup scripts)

### Setup

1. **Generate SSL certificate and start services:**
   ```bash
   bash setup.sh
   ```

2. **Add domain to hosts file:**
   Add this line to your hosts file (`/etc/hosts` on Linux/Mac, `C:\Windows\System32\drivers\etc\hosts` on Windows):
   ```
   127.0.0.1 hospital-a.api.co.th
   ```

3. **Access the application:**
   - HTTP: http://hospital-a.api.co.th (redirects to HTTPS)
   - HTTPS: https://hospital-a.api.co.th
   - Direct Go app: http://localhost:8080

### Manual Setup

If you prefer to run commands manually:

1. **Generate SSL certificate:**
   ```bash
   mkdir -p ssl logs/nginx
   bash scripts/generate-ssl.sh
   ```

2. **Start services:**
   ```bash
   docker-compose up -d
   ```

## Services

- **agnos_app**: Go application (port 8080)
- **agnos_nginx**: Nginx reverse proxy (ports 80, 443)
- **agnos_db**: PostgreSQL database (port 5432)

## SSL Certificate

The setup uses a self-signed certificate for `hospital-a.api.co.th`. Your browser will show a security warning - this is normal for self-signed certificates. Click "Advanced" and "Proceed to hospital-a.api.co.th" to continue.

## Development

To rebuild the Go application:
```bash
docker-compose build app
docker-compose up -d app
```

## Logs

View logs for specific services:
```bash
docker-compose logs app
docker-compose logs nginx
docker-compose logs db
```

## Stopping Services

```bash
docker-compose down
```