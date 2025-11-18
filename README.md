# Local RAG API

A Go-based API service for indexing and searching repositories using Retrieval-Augmented Generation (RAG) techniques.

Based on [LangChain framework](https://langchain.com) principles.

![rag_detail](https://github.com/langchain-ai/rag-from-scratch/assets/122662504/54a2d76c-b07e-49e7-b4ce-fc45667360a1)

## Roadmap

✅ Health check endpoint  
✅ GitHub repository indexing  
⬜ Discussion management (CRUD operations)  
⬜ Question/answer within discussions  
⬜ Discussion history tracking  

## Features

- 🔍 **GitHub Repository Indexing**: Fetch and index content from public GitHub repositories
- 📄 **File Type Filtering**: Selectively index specific file types (e.g., `.md`, `.mdx`)
- 🚀 **REST API**: Simple HTTP endpoints for indexing and health checks
- 📚 **Swagger Documentation**: Auto-generated API documentation
- 🧱 **DDD Architecture**: Domain-driven design with clear separation of concerns
- 🔌 **Vector Store Integration**: ChromaDB for embedding storage
- 🤖 **Ollama Integration**: Local embeddings generation

## Prerequisites

- Go 1.23.3 or higher.  
- Git (for cloning repositories).  
- [swag](https://github.com/swaggo/swag) CLI tool for Swagger docs.  
- [Ollama](https://ollama.ai/) for local embeddings.  
- [ChromaDB](https://www.trychroma.com/) for vector storage.  
- [SQLite](https://www.sqlite.org/index.html) for local database storage.  

Install swag:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### DDD Architecture Layers
- **`cmd/`**: Application entry points
- **`internal/api/`**: Interface layer - HTTP handlers, routing, request/response types
- **`internal/application/`**: Application layer - Use cases coordinating domain logic
- **`internal/domain/`**: Domain layer - Business logic (embeddings, fetcher, indexer)
- **`internal/infra/`**: Infrastructure layer - External services and persistence

## Quick Start

### 1. Clone the Repository

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
make run
```

The server will start on `http://localhost:8080`

### Build Binary

To create a production binary:

```bash
make build
```

The binary will be available at `bin/server`

## API Endpoints

### Health Check
```http
GET /api/v1/health
```

### Index GitHub Repository
```http
POST /api/v1/index-github
Content-Type: application/json

{
  "repository": "owner/repo",
  "branch": "main"
}
```

### Swagger Documentation
```http
GET /swagger/index.html
```

## Development

### Run Tests
```bash
make test
```

### Format Code
```bash
make fmt
```

### Lint Code (requires golangci-lint)
```bash
make lint
```

### Clean Build Artifacts
```bash
make clean
```

## Configuration

Default configuration:

- **Server Port**: `8080`
- **API Base Path**: `/api/v1`
- **Swagger UI**: `/swagger/index.html`

## License

MIT License

