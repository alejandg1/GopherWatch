# GopherWatch

GopherWatch is a lightweight container monitoring tool written in Go. It provides real-time statistics for Docker containers, including CPU and memory usage. It exposes these metrics via a simple HTTP API, making it easy to integrate with various status bars and dashboards.

## Features

- Real-time container monitoring.
- CPU and Memory usage statistics.
- Lightweight and efficient.
- Simple HTTP API for data consumption.
- Designed for integration with window manager status bars (e.g., Waybar, Eww).

## Requirements

- Docker
- Go 1.21 or higher (for building from source)

## Installation

### Building from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/actiometa/gopherwatch.git
   cd gopherwatch
   ```

2. Build the application:
   ```bash
   go build -o gopherwatch cmd/gopherwatch/main.go
   ```

## Usage

1. Ensure Docker is running.

2. Run the application:
   ```bash
   ./gopherwatch
   ```

   The server will start on port 8080 by default.

## API Documentation

### GET /v1/stats

Returns a JSON array containing statistics for all running containers.

**Response Format:**

```json
[
  {
    "id": "container_id",
    "name": "container_name",
    "status": "running",
    "cpu_percentage": 0.5,
    "mem_usage_mb": 128.5,
    "mem_limit_mb": 1024.0,
    "mem_percentage": 12.5
  }
]
```
