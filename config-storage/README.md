# Config Storage Service

A configuration storage service that supports MinIO bucket storage for distributed configuration management.

## Features

- RESTful API for CRUD operations on configurations
- MinIO bucket storage backend (default)
- File-based storage fallback option
- JSON serialization support
- Concurrent-safe operations

## Building

```bash
go mod download
go build -o config-storage ./cmd
```

## Running

### With MinIO Storage (Default)

Set environment variables:

```bash
export STORAGE_TYPE=minio
export MINIO_ENDPOINT=localhost:9000
export MINIO_ACCESS_KEY=minioadmin
export MINIO_SECRET_KEY=minioadmin
export MINIO_BUCKET=configs
export MINIO_USE_SSL=false
```

Then run:

```bash
go run ./cmd/main.go
```

### With File Storage

Set environment variables:

```bash
export STORAGE_TYPE=file
export DATA_DIR=./data
```

Then run:

```bash
go run ./cmd/main.go
```

## API Endpoints

- `POST /configs` - Create a new configuration
- `GET /configs` - List all configurations
- `GET /configs/{id}` - Get a specific configuration
- `PUT /configs/{id}` - Update a configuration
- `DELETE /configs/{id}` - Delete a configuration

## Configuration Model

```json
{
  "id": "config-1",
  "name": "My Config",
  "type": "json",
  "environment": "production",
  "json_content": { ... },
  "tags": ["tag1", "tag2"]
}
```

## Docker Compose Example

To run with MinIO locally:

```yaml
version: '3.8'
services:
  minio:
    image: minio/minio:latest
    ports:
      - "9000:9000"
      - "9001:9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    command: server /data --console-address ":9001"
```

## Dependencies

- github.com/gorilla/mux - HTTP router
- github.com/minio/minio-go/v7 - MinIO client SDK
