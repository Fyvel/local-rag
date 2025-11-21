# Copilot Instructions — Implementing Discussion Routes (Go DDD)

## Goal
Implement the upcoming "discussion" feature in the existing Go DDD architecture.  
The feature includes routes for:
- Listing discussions
- Creating a discussion
- Fetching a discussion by ID
- Asking a question in a discussion
- Retrieving discussion history

The routes currently appear commented in `api/router.go`.  
Copilot must implement them **end-to-end** across layers:

- `api/` handlers (Gin)
- `application/` use cases
- `domain/` models, aggregates, factories, and repositories
- `infra/` repository implementations
- `di/` wiring (constructors + initialization)

Every step must compile, be incremental, and preserve clean DDD boundaries.

## Architectural Expectations

### Domain layer (`domain/discussion`)
Copilot must:
- Introduce a `Discussion` aggregate (ID, title/topic, message list, timestamps).
- Introduce a `Message` or `Entry` value object as needed.
- Implement domain invariants such as:
  - Cannot add empty messages.
  - Cannot add messages to a closed discussion (if applicable).
- Define repository interfaces (e.g., `DiscussionRepository`).
- Keep domain models free of infrastructure or API concerns.

### Application layer (`application/discussion`)
Copilot must:
- Implement use cases for:
  - `ListDiscussions`
  - `CreateDiscussion`
  - `GetDiscussion`
  - `AskQuestion`
  - `GetDiscussionHistory`
- Inject domain repository interface(s).
- Translate DTOs from/to domain objects.
- Handle orchestration logic cleanly.

### API layer (`api/`)
Copilot must:
- Introduce handler factories similar to existing patterns.
- Decode requests into application DTOs.
- Encode results into JSON responses.
- Implement or improve the following endpoints:

```
GET /api/v1/discussions
POST /api/v1/discussions
GET /api/v1/discussions/:id
POST /api/v1/discussions/:id/question
GET /api/v1/discussions/:id/history
```

- Ensure handlers never contain domain or infrastructure logic.

### Infra layer (`infra/discussion`)
Copilot must:
- Create concrete repository implementations.
- Provide minimal working storage (in-memory first).
- Allow later extension to database or AI vector store.
- Keep JSON, DB, filesystem, or external service logic here.

### DI layer (`di/`)
Copilot must:
- Add factories to build all discussion use cases.
- Initialize domain repo via infra implementation.
- Provide the instance to API factories.

## Route Implementation Rules

When implementing or adjusting the routes:
- Ensure consistent naming (Go conventions).
- Prefer early exits in handlers.
- Validate path parameters and request bodies at handler level.
- Use consistent response shapes:
- `{ "data": … }` or `{ "error": … }`
- Avoid leaking infra types or domain types directly to HTTP responses.

Copilot may modify the commented endpoints if an improved naming pattern is clearer:
- `AskQuestionHandler` → `CreateMessageHandler` (if better domain alignment)
- `HistoryHandler` → `ListMessagesHandler` (if aligned to domain naming)

## Iteration Requirements

For each proposed change, Copilot must:
1. Provide a short explanation describing the architectural reasoning.
2. Show complete updated files or new files.
3. List movements, renames, or deletions.
4. Generate a suggested commit message.
5. Provide brief trade-off notes when a design decision is not obvious.

## Additional Guidance

- Start with domain models before application or API code.
- Prefer small commits per layer or per use case.
- Avoid exposing complex domain internals in HTTP responses.
- Keep error handling minimal and explicit.
- Ensure `NewRouter` stays clean, delegating logic to factories.

## Acceptable Enhancements
Copilot is allowed to:
- Add request/response DTOs.
- Introduce pagination for list endpoints.
- Introduce domain events (optional).
- Suggest better naming for endpoints, handlers, and use cases.
- Add lightweight middleware if needed (e.g., request validation).

## Completion Criteria
The feature is considered fully implemented when:
- All discussion routes compile and function end-to-end.
- Domain models reflect the discussion domain accurately.
- Use cases encapsulate logic with no leakage.
- Infrastructure has minimal working repository implementation.
- DI container can construct all required services.
- The router exposes all routes cleanly.
- No layer violates DDD boundaries.

