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

### 🎯 Recommended Next Steps
The following steps maintain incremental progress toward DDD architecture:

#### Next: Step 2 - Extract Embedding value object to domain
- Move `Embedding` struct from `embeddings` package to `domain/embeddings`
- Move Base64 decoder logic as a domain method
- Keep `Embedder[T]` interface in embeddings (it's infrastructure-facing)
- Update references in `ollama` package and others

#### Future: Step 3 - Extract Repository domain entities
- Move `GithubRepository` and `GithubFile` from `fetcher` to `domain/repositories`
- These are domain concepts, not infrastructure
- Keep git/filesystem operations in `fetcher` as infrastructure

#### Future: Step 4 - Define repository interfaces in domain
- Create `domain/documents/repository.go` with DocumentRepository interface
- Define methods: Add, Remove, Query, etc.
- Move interface to domain, keep implementation in `store`

#### Future: Step 5 - Move concrete store implementations to infra
- Reorganize `internal/store/chroma` → `internal/infra/store/chroma`
- Keep VectorStore in `store` but aligned with domain interface

#### Future: Step 6 - Introduce application use cases
- Create `application/indexing/index_repository.go` use case
- Extract orchestration logic from `indexer.IndexGithubRepository`
- Handler calls application use case instead of indexer directly

#### Future: Step 7 - Domain services for business logic
- Extract validation, chunking logic into domain services
- Move embedding coordination into domain service

Each step should end with compiling code and an example commit message.

## Output Formatting
Copilot responses should include:
- A short explanation of the proposed incremental change.
- Updated or newly created files with complete code blocks.
- A suggested commit message.
- Notes about trade-offs or impacts when relevant.

