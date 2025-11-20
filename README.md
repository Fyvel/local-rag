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
- 🔌 **Vector Store Integration**: ChromaDB for embedding storage
- 🤖 **Ollama Integration**: Local embeddings generation

## Prerequisites

- Go 1.23.3 or higher  
- Git (for cloning repositories)  
- [swag](https://github.com/swaggo/swag) CLI tool for Swagger docs  
- [Ollama](https://ollama.ai/) for local embeddings  
- [ChromaDB](https://www.trychroma.com/) for vector storage  

Install swag:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## Architecture

This project follows Domain-Driven Design (DDD) principles with clear layer separation and proper dependency flow.

### Layer Overview

```
┌─────────────────────────────────────────────────────────┐
│                    server/ (main.go)                    │
│                     Entry Point                         │
└─────────┬───────────────────────────────────────────────┘
          │
          ↓
┌─────────────────────────────────────────────────────────┐
│                    internal/di/                         │
│              Dependency Injection                       │
│         Wires up all application components             │
└─────────┬───────────────────────────────────────────────┘
          │
          ↓
┌─────────────────────────────────────────────────────────┐
│                   internal/api/                         │
│     Interface Layer (HTTP handlers & routing)           │
│   • Handles HTTP requests/responses                     │
│   • Validates input                                     │
│   • Delegates to application layer                      │
│   Dependencies: application                             │
└─────────┬───────────────────────────────────────────────┘
          │
          ↓
┌─────────────────────────────────────────────────────────┐
│               internal/application/                     │
│          Application Layer (Use Cases)                  │
│   • Orchestrates domain logic                           │
│   • Coordinates workflows                               │
│   • Transaction management                              │
│   Dependencies: domain only                             │
└─────────┬───────────────────────────────────────────────┘
          │
          ↓
┌─────────────────────────────────────────────────────────┐
│                  internal/domain/                       │
│        Domain Layer (Business Logic)                    │
│   • Core business entities                              │
│   • Domain services                                     │
│   • Repository interfaces                               │
│   • Business rules & invariants                         │
│   Dependencies: none (pure domain)                      │
└─────────────────────────────────────────────────────────┘
          ↑
          │ (implements)
          │
┌─────────┴───────────────────────────────────────────────┐
│                   internal/infra/                       │
│       Infrastructure Layer (External Systems)           │
│   • Database implementations                            │
│   • External API clients                                │
│   • File system operations                              │
│   • Third-party integrations                            │
│   Dependencies: domain (implements interfaces)          │
└─────────────────────────────────────────────────────────┘
```

### Layer Responsibilities

#### `server/`
Application entry point. Initializes the DI container with configuration and starts the HTTP server.

#### `internal/api/` - Interface Layer
- HTTP handlers using Gin framework
- Request/response DTOs
- Input validation
- HTTP routing
- **Can depend on:** `application/` only
- **Cannot depend on:** `domain/` directly, `infra/`, `di/`

#### `internal/application/` - Application Layer
- Use cases (e.g., `IndexRepositoryUseCase`)
- Orchestration logic
- Transaction boundaries
- Workflow coordination
- **Can depend on:** `domain/` interfaces only
- **Cannot depend on:** `api/`, `infra/`, `di/`

#### `internal/domain/` - Domain Layer
Core business logic with four main aggregates:
- **`documents/`**: Document entities, repository interface, domain service, query results
- **`repositories/`**: Repository aggregate, fetcher interface
- **`embeddings/`**: Embedding value objects, embedder interface
- **`transformers/`**: Text chunking interface (abstraction for text splitting logic)
- **Cannot depend on:** any other layer (pure domain)

#### `internal/infra/` - Infrastructure Layer
Concrete implementations of domain interfaces:
- **`embeddings/ollama/`**: Ollama embeddings client (implements `embeddings.Embedder`)
- **`fetcher/`**: GitHub repository fetcher (implements `repositories.Fetcher`)
- **`store/chroma/`**: ChromaDB vector store client
- **`store/vectorstore/`**: Vector store adapter (implements `documents.DocumentRepository`)
- **`textsplitter/`**: Text splitting adapter wrapping langchaingo (implements `textsplitter.TextSplitter`)
- **`httpclient/`**: HTTP client utilities
- **Can depend on:** `domain/` (to implement interfaces)
- **Cannot depend on:** `api/`, `application/`

#### `internal/di/` - Dependency Injection
Container that wires up all dependencies. Can import everything but contains no business logic.

### Dependency Rules

The project strictly enforces proper dependency flow:

✅ **Allowed:**
- `api` → `application`
- `application` → `domain` (interfaces only)
- `infra` → `domain` (implements interfaces)
- `di` → all layers (for wiring only)

❌ **Forbidden:**
- `domain` → anything (must be pure, no internal imports)
- `application` → `infra` (must use domain interfaces instead)
- `application` → `api`
- `api` → `domain` (must go through application layer)
- `api` → `infra` (except in main/di setup)

### DDD Compliance

This project follows Domain-Driven Design principles:

1. **Pure Domain Layer**: The domain layer has zero dependencies on other internal packages. All domain logic is expressed through entities, value objects, and interfaces.

2. **Dependency Inversion**: Infrastructure implements domain interfaces, not the other way around. The domain defines contracts (`Fetcher`, `Embedder`, `DocumentRepository`) that infrastructure fulfills.

3. **Application Services as Use Cases**: Each use case (e.g., `IndexRepositoryUseCase`) orchestrates domain operations without knowing implementation details.

4. **Hexagonal Architecture**: The domain is at the center, surrounded by ports (interfaces) and adapters (infrastructure implementations).

5. **Explicit Context Boundaries**: Clear package boundaries with `documents/`, `embeddings/`, and `repositories/` aggregates.

### Key Design Patterns

1. **Dependency Inversion**: All cross-layer dependencies point inward toward the domain
2. **Repository Pattern**: Data access abstracted through domain interfaces
3. **Use Case Pattern**: Application services orchestrate domain logic
4. **Factory Pattern**: Handler factories receive injected dependencies
5. **Adapter Pattern**: Infrastructure adapters implement domain interfaces

### Architectural Improvements

This project has undergone rigorous architectural review to ensure full DDD compliance:

#### ✅ **Clean Architecture Enforcement**
- **Zero external dependencies in domain layer**: The domain layer has no imports from `application`, `api`, `infra`, or `di`
- **Application layer independence**: Application logic depends only on domain interfaces, never concrete infrastructure
- **Pure domain interfaces**: All infrastructure concerns abstracted behind domain-defined contracts

#### ✅ **Type Safety**
- **No `interface{}` returns**: Repository methods return properly typed domain models
- **Query results use domain types**: `QueryResponse` and `QueryResult` structs ensure type safety
- **Strong typing throughout**: All interfaces and implementations use concrete types where possible

#### ✅ **Infrastructure Abstraction**
- **TextChunker interface**: Extracted third-party text splitting dependency into domain `transformers` package
- **Adapter pattern for external libraries**: Infrastructure adapters wrap external libraries (langchaingo textsplitter, etc.)
- **Swappable implementations**: Any infrastructure component can be replaced without touching domain or application layers

#### ✅ **No Circular Dependencies**
- Project builds cleanly with `go build ./...`
- Clear unidirectional dependency flow
- All import cycles prevented by design

### Go Package Conventions

This project follows Go best practices for package organization:

1. **Interfaces in consumer packages**: Domain interfaces are defined where they're needed, not where they're implemented
2. **Small, focused interfaces**: Following Go's "accept interfaces, return structs" principle
3. **Package naming**: Packages use descriptive names (`documents`, `embeddings`, `repositories`)
4. **Avoid package stuttering**: Types are named `Document` not `DocumentDocument`
5. **Clear dependency direction**: Dependencies always point toward stability (domain is most stable)

### Package Import Rules

```go
// ✅ GOOD: Application depends on domain interfaces
import "local-ai/internal/domain/documents"
import "local-ai/internal/domain/embeddings"

// ❌ BAD: Domain should never import from other layers
import "local-ai/internal/application/..."  // FORBIDDEN in domain/
import "local-ai/internal/infra/..."        // FORBIDDEN in domain/

// ✅ GOOD: Infrastructure implements domain interfaces
import "local-ai/internal/domain/embeddings"

// ❌ BAD: Infrastructure shouldn't know about application
import "local-ai/internal/application/..."  // FORBIDDEN in infra/
```

## Quick Start

### 1. Prerequisites Setup

Start required services (ChromaDB and Ollama):

```bash
# Start ChromaDB (using Docker)
docker run -p 8000:8000 chromadb/chroma

# Start Ollama and pull the embedding model
ollama serve
ollama pull mxbai-embed-large
```

### 2. Clone and Build

```bash
git clone https://github.com/fyvel/local-rag.git
cd local-rag

# Install dependencies
make deps

# Generate Swagger docs
swag init -g server/main.go -o docs

# Build the server
make build
```

### 3. Run the Server

```bash
make run
```

The server will start on `http://localhost:8080`

Alternatively, run the binary directly:

```bash
./bin/server
```

## API Endpoints

### Health Check
```http
GET /api/v1/health
```

Response:
```json
{
  "success": true,
  "data": {
    "status": "ok",
    "service": "local-rag-api",
    "version": "1.0.0"
  }
}
```

### Index GitHub Repository
```http
POST /api/v1/index-github
Content-Type: application/json

{
  "github_url": "https://github.com/username/repo",
  "file_types": ["md", "mdx"],
  "index_source": "my_index"
}
```

Response:
```json
{
  "success": true,
  "data": {
    "index_id": "my_index"
  }
}
```

### Swagger Documentation
Interactive API documentation available at:
```
http://localhost:8080/swagger/index.html
```

## Configuration

Configuration is set in `server/main.go` via the DI container. Default values:

| Setting | Default Value | Description |
|---------|---------------|-------------|
| Server Port | `8080` | HTTP server port |
| ChromaDB URL | `http://localhost:8000` | Vector database endpoint |
| Ollama URL | `http://localhost:11434/api` | Embeddings service endpoint |
| Target Index | `default_index` | Default collection name |
| Embedding Model | `mxbai-embed-large` | Model for generating embeddings |

To customize, edit the `di.Config` in `server/main.go`:

```go
cfg := di.Config{
    ChromaURL:   os.Getenv("CHROMA_URL"),
    OllamaURL:   os.Getenv("OLLAMA_URL"),
    TargetIndex: os.Getenv("DEFAULT_INDEX"),
}
```

## Development

### Project Structure
```
.
├── server/                 # Application entry point
│   └── main.go
├── internal/
│   ├── api/               # HTTP handlers & routing
│   ├── application/       # Use cases
│   ├── domain/            # Business logic
│   ├── infra/             # External integrations
│   └── di/                # Dependency injection
├── docs/                  # Generated Swagger docs
└── Makefile
```

### Development Commands

```bash
# Run tests
make test

# Format code
make fmt

# Lint code (requires golangci-lint)
make lint

# Clean build artifacts
make clean

# Regenerate Swagger docs
swag init -g server/main.go -o docs
```

### Adding New Features

When extending the application, follow these guidelines:

1. **New Domain Concepts**: Add to `internal/domain/`
   - Define entities and value objects
   - Create repository interfaces
   - Implement domain services

2. **New Use Cases**: Add to `internal/application/`
   - Create use case structs
   - Depend only on domain interfaces
   - Keep orchestration logic here

3. **New Infrastructure**: Add to `internal/infra/`
   - Implement domain interfaces
   - Handle external system integration
   - Keep implementation details isolated

4. **New Endpoints**: Add to `internal/api/`
   - Create handler factories
   - Define request/response types
   - Add routes to router.go

5. **Wire Dependencies**: Update `internal/di/container.go`
   - Register new implementations
   - Create factory methods
   - Maintain single source of dependency configuration

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./internal/domain/documents/...
```

## Extending the Architecture

### Adding a New Repository Source

To support sources beyond GitHub (e.g., GitLab, Bitbucket):

1. Implement `repositories.Fetcher` interface in `internal/infra/fetcher/`
2. Register in DI container
3. No changes needed in application or domain layers

### Adding a New Embedding Provider

To use different embedding services (e.g., OpenAI, HuggingFace):

1. Implement `embeddings.Embedder` interface in `internal/infra/embeddings/`
2. Register in DI container
3. Application layer remains unchanged

### Adding a New Vector Store

To use different vector databases (e.g., Pinecone, Weaviate):

1. Implement `documents.DocumentRepository` interface in `internal/infra/store/`
2. Register in DI container
3. Domain and application layers remain unchanged

## Troubleshooting

### ChromaDB Connection Error
```
failed to ensure collection: connection refused
```
**Solution**: Ensure ChromaDB is running on port 8000

### Ollama Embedding Error
```
embedding failed: model not found
```
**Solution**: Pull the required model: `ollama pull mxbai-embed-large`

### GitHub Clone Error
```
failed to clone repository: authentication required
```
**Solution**: Use public repositories or configure Git credentials

## License

MIT License

