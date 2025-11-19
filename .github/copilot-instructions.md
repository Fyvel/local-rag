# Copilot Refactoring Instructions (Go + DDD Migration)

## Goal
Incrementally refactor the Go codebase into a clean Domain-Driven Design (DDD) architecture using the final folder structure:

- `api/` — HTTP, gRPC, CLI, or any delivery layer.
- `application/` — use cases, orchestrators, DTOs, application services.
- `domain/` — entities, value objects, aggregates, domain services, domain errors, repository interfaces.
- `infra/` — concrete implementations: database adapters, external clients, messaging, logging.
- `di/` — composition root: wiring dependencies (manual or with a DI tool).

The end state must:
- Keep the project compiling at every step.
- Preserve behavior unless intentionally modified.
- Follow Go idioms: small packages, clear naming, minimal interfaces, explicit dependencies.
- Fully relocate code to the structure above and remove obsolete folders.

## Architectural Expectations
**Domain**
- Must be framework-agnostic.
- Must not import `infra`, `api`, or `application`.
- Define repository and service interfaces only when required by domain rules.
- Model aggregates carefully; keep invariants inside domain constructors and methods.

**Application**
- Calls domain services and repositories.
- Contains use cases, commands, queries, orchestrators.
- Translates between API DTOs and domain objects.

**Infra**
- Concrete implementations for repositories, external services, and technical details.
- Keep infra-specific errors wrapped and translated into domain/application errors at boundaries.

**API**
- Only delivery concerns: HTTP routing, JSON encoding/decoding, gRPC definitions, CLI commands.
- No business logic.

**DI**
- Central initialization: constructing services, repositories, handlers, application use cases.
- No domain logic here.

## Refactoring Rules
To propose or apply changes, Copilot should:

1. **Work incrementally**
   - One clear boundary improvement per iteration.
   - Produce compiling, runnable code at each step.
   - Avoid speculative abstractions.

2. **Prioritize structure clarity**
   - Move code into `domain` first (entities, logic).
   - Introduce `application` when orchestrations emerge.
   - Create `infra` when relocating DB or external integrations.
   - Move handlers/controllers into `api`.
   - Centralize wiring into `di`.

3. **Finalization phase**
   Copilot must:
   - Identify any remaining files outside the target folders.
   - Relocate them to the appropriate layer.
   - Remove old or duplicate folders.
   - Rename packages to follow Go conventions (lowercase, no plural unless meaningful, cohesive packages).
   - Update imports across the project.
   - Ensure no circular dependencies remain.
   - Clean up any leftover transitional code or adapters.

4. **Go conventions**
   - Keep package names concise and domain-oriented.
   - Avoid overly deep package hierarchies.
   - Avoid stuttered names (e.g., `order.OrderService` is preferred over `order.OrderDomainService`).
   - Prefer small, explicit interfaces placed where they are consumed.

5. **Preserve behavior**
   - Refactor structure, not semantics.
   - Maintain handler signatures unless intentionally adjusting APIs.
   - Keep tests running; update them as directories move.

## Code Generation Expectations
Generated or updated code must:
- Be complete, not pseudocode.
- Include minimal working implementations when required (no broken stubs).
- Use clear constructors for services and aggregates.
- Avoid unnecessary generics or abstractions.
- Respect runtime efficiency and error clarity.

## Step-by-Step Iteration Examples
Examples of valid steps:
- Extract an entity or value object into `domain/<context>`.
- Move business logic from handlers/controllers to a domain service.
- Introduce a repository interface in `domain`, then move DB code to `infra`.
- Create an application use case and update the API handler to call it.
- Introduce a `di` package and consolidate object construction.
- Relocate leftover utility packages to their correct layers or remove them.
- Normalize naming and imports after full folder relocation.

Each step should end with compiling code and a suggested commit message.

## Output Format
Copilot responses must include:
1. A short explanation of the change.
2. Complete code blocks for modified or new files.
3. Files that must be deleted or moved.
4. A suggested commit message.
5. Notes on trade-offs when relevant.

