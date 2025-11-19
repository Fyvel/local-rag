# Copilot Refactoring Instructions (Go + DDD Migration)

## Goal
Incrementally refactor the existing Go codebase toward a Domain-Driven Design (DDD) architecture.  
Each refactoring step must:
- Produce compiling, runnable code.
- Be small enough to commit independently.
- Preserve existing behavior unless explicitly modifying it.
- Improve structure, boundaries, and domain clarity.

## Architectural Direction
The target architecture uses:
- **Api**: REST API handlers, routing, request/response types.
- **Domain**: Entities, value objects, domain services, domain errors, aggregates.
- **Application**: Use cases, orchestrating domain logic, DTOs.
- **Infrastructure**: Repositories, persistence adapters, external services, frameworks.

Use package naming aligned with features/context (e.g., `chat`, `documents`, `embeddings`, `tools`), not technical layers alone.
Avoid circular dependencies; domain must not depend on infrastructure.

## Refactoring Rules
When suggesting code changes, Copilot should:
1. **Propose small, incremental steps**
   - One clear responsibility per commit.
   - Example commit types: extract domain model, isolate repository interface, move business logic into domain service, introduce application command/query, etc.

2. **Ensure the build stays green**
   - No broken imports.
   - No placeholder stubs that break compilation.  
   - If a dependency is not implemented yet, provide a minimal working implementation.

3. **Protect existing behavior**
   - Maintain API signatures unless refactoring requires otherwise.
   - Migrate logic without altering semantics.

4. **Promote clean boundaries**
   - Push business rules into the domain layer.
   - Move technical code out of the domain.
   - Introduce interfaces only where they provide clear value.

5. **Encourage Go best practices**
   - Keep packages small and composable.
   - Avoid unnecessary abstractions.
   - Favor simple constructors and explicit dependencies.
   - Optimize for readability and runtime efficiency.

## Code Generation Guidelines
When generating or rewriting Go code:
- Prefer pure functions and immutable value objects when feasible.
- Keep dependencies injected, not global.
- Avoid large god-structs; keep aggregates tight.
- Avoid leaking infrastructure concerns into domain types.
- Handle errors explicitly and clearly.

## Refactoring Progress

### ✅ Completed Steps
1. **Extract Document entity to domain layer** (commit: `refactor: extract Document entity into domain layer`)
   - Created `internal/domain/documents/document.go`
   - Moved Document struct from `store` package to domain
   - Moved ID generation functions (GenerateID, GenerateDeterministicID) to domain
   - Updated `store.VectorStore` and `indexer` to use `domain.Document`
   - Status: ✅ Compiles, no behavior changes

2. **Extract Embedding value object to domain layer** (commit: `refactor: extract Embedding value object into domain layer`)
   - Created `internal/domain/embeddings/embedding.go`
   - Moved Embedding struct from `embeddings` package to domain
   - Moved Base64String type and Decode logic to domain
   - Added NewEmbedding constructor with validation and immutability
   - Kept `Embedder[T]` interface in embeddings package (infrastructure-facing)
   - Updated `ollama` package to use `domain/embeddings.Embedding`
   - Updated tests to use domain types
   - Status: ✅ Compiles, tests pass, no behavior changes

3. **Extract Repository domain entities** (commit: `refactor: extract Repository domain entities into domain layer`)
   - Created `internal/domain/repositories/repository.go`
   - Moved `GithubRepository` → `Repository` and `GithubFile` → `File` to domain
   - Changed `Repository` to be an aggregate root containing `File` entities
   - Added domain methods: `FileCount()`, `HasFiles()`
   - Updated `fetcher.GetGithubRepository` to return `*repositories.Repository`
   - Updated `indexer` to work with domain repository types
   - Kept git/filesystem operations in `fetcher` (infrastructure adapter)
   - Status: ✅ Compiles, no behavior changes

4. **Define repository interfaces in domain** (commit: `refactor: define DocumentRepository interface in domain`)
   - Created `internal/domain/documents/repository.go` with DocumentRepository interface
   - Defined interface methods: AddDocument, RemoveDocument, Query
   - Updated `store.VectorStore` to implement DocumentRepository interface
   - Modified RemoveDocument to accept context parameter
   - Modified Query to return interface{} for infrastructure flexibility
   - Updated `indexer.Indexer` to depend on DocumentRepository interface instead of concrete VectorStore
   - Applied dependency inversion principle (domain defines interface, infrastructure implements)
   - Status: ✅ Compiles, no behavior changes

5. **Move concrete store implementations to infra** (commit: `refactor: move ChromaDB implementation to infrastructure layer`)
   - Moved `internal/store/chroma` → `internal/infra/store/chroma`
   - Updated all imports from `store/chroma` to `infra/store/chroma`
   - ChromaDB client now clearly organized as infrastructure adapter
   - VectorStore remains in `store` package as adapter implementing domain interface
   - Separated infrastructure concerns from domain/application layers
   - Status: ✅ Compiles, no behavior changes

6. **Introduce application use cases** (commit: `refactor: introduce application layer with IndexRepositoryUseCase`)
   - Created `internal/application/indexing/index_repository.go`
   - Introduced `IndexRepositoryUseCase` to orchestrate repository indexing workflow
   - Defined `IndexRepositoryCommand` and `IndexRepositoryResult` DTOs
   - Introduced `RepositoryFetcher` interface for dependency abstraction
   - Exported `indexer.ProcessRepository` method for use by application layer
   - Updated API handler to use application use case instead of calling indexer directly
   - Separated orchestration (application) from low-level processing (indexer)
   - Status: ✅ Compiles, no behavior changes

### 🎯 Recommended Next Steps
The following steps maintain incremental progress toward DDD architecture:

#### Next: Step 7 - Domain services for business logic
- Extract validation, chunking logic into domain services
- Move embedding coordination into domain service

Each step should end with compiling code and an example commit message.

## Output Formatting
Copilot responses should include:
- A short explanation of the proposed incremental change.
- Updated or newly created files with complete code blocks.
- A suggested commit message.
- Notes about trade-offs or impacts when relevant.

