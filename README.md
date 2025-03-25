# Face Recognition Processor

A Go-based facial recognition system with REST API endpoints for training and comparing faces.

## Features

- Face training with unique identifiers
- Face comparison with similarity scores
- RESTful API endpoints
- Base64 image support
- Local storage of face embeddings
- Docker support

## Requirements

- Docker and Docker Compose
- Go 1.21+ (for local development)

## Quick Start

1. Clone the repository
2. Create a `models` directory and download the required face recognition models
3. Run with Docker Compose:

```bash
docker-compose up --build
```

The API will be available at `http://localhost:8080`

## API Endpoints

### Train Face

```http
POST /train
Content-Type: application/json

{
    "image": "base64_encoded_image",
    "person_id": "unique_identifier"
}
```

### Compare Face

```http
POST /compare
Content-Type: application/json

{
    "image": "base64_encoded_image"
}
```

Response:
```json
{
    "person_id": "matched_identifier",
    "similarity": 0.95
}
```

## Development

1. Install dependencies:
```bash
go mod download
```

2. Run locally:
```bash
go run cmd/api/main.go
```

## Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── api/
│   │   └── handler.go
│   ├── face/
│   │   └── service.go
│   └── storage/
│       └── service.go
├── models/           # Face recognition models
├── data/            # Stored face embeddings
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## Configuration

The application uses the following default configuration:

- Server port: 8080
- Face similarity threshold: 0.6
- Storage location: ./data/faces.json

## License

MIT License
