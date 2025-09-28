# Local RAG API

A Go-based API service for indexing and searching GitHub repositories using Retrieval-Augmented Generation techniques.

## Overview

This project provides a REST API that can fetch and index GitHub repositories, making their contents searchable for RAG applications. The service extracts content from repositories and processes it for embedding-based search.

## Features

- 🔍 **GitHub Repository Indexing**: Fetch and index content from public GitHub repositories
	- 📄 **File Type Filtering**: Selectively index specific file types (e.g., `.md`, `.mdx`)
- 🚀 **REST API**: Simple HTTP endpoints for indexing operations
- 📚 **Swagger Documentation**: Auto-generated API documentation

## Prerequisites

- Go 1.23.3 or higher
- Git (for cloning repositories)
- [swag](https://github.com/swaggo/swag) CLI tool for generating Swagger docs

Install swag:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## Project Structure

```
src/
├── cmd/server/          # Application entry point
├── api/                 # API handlers and types
├── internal/           # Internal packages
│   ├── client/         # HTTP client utilities
│   ├── embeddings/     # Embedding generation
│   ├── fetcher/        # Repository fetching logic
│   ├── indexer/        # Content indexing logic
│   ├── store/          # Data storage (Chroma)
│   └── ...
├── docs/               # Auto-generated Swagger documentation
```

## Quick Start

### 1. Clone and Navigate

```bash
git clone https://github.com/fyvel/local-rag.git
cd local-rag/src
```

### 2. Install Dependencies

```bash
make deps
```

### 3. Build and Run

```bash
# Build the project
make build

# Or run directly (includes docs generation)
make run

# Or use the development workflow
make dev
```

### 4. Alternative: Run without Make

```bash
# Generate Swagger docs
swag init -g cmd/server/main.go -o docs

# Run the server
go run ./cmd/server
```

The server will start on `http://localhost:8080`


## API Documentation
```http
GET /swagger/index.html
```

## Configuration

The application currently uses default configuration. Key settings:

- **Server Port**: 8080
- **API Base Path**: `/api/v1`
- **Swagger UI**: `/swagger/index.html`

