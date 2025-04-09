# Osquery Data Updater

This project periodically updates the system's OS version, osquery version, and installed apps, and exposes an API endpoint (`/latest_data`) to retrieve the latest data.

## Prerequisites

Before running this project, make sure you have the following installed:
- Go (Golang)
- osquery
- Docker

## Getting Started

1. **Start the Docker Container:**
   Run the following command to bring up the PostgreSQL container using Docker Compose:
   ```bash
   docker-compose up

2. **Start the Gin Server:**
   Once the container is running, you can start the Gin server by executing:
    ```bash
    go run cmd/server/main.go