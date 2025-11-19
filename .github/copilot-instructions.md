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

## Step-by-Step Iteration Examples
Copilot can propose iterations such as:
- Step 1: Introduce a `domain/<context>` package and move core structs there.
- Step 2: Extract business logic from handlers into domain services.
- Step 3: Define repository interfaces in the domain.
- Step 4: Move concrete DB code into `infra/<adapter>`.
- Step 5: Add application use cases in `application/<context>`.
- Step 6: Update handlers to call application services instead of domain/internal logic.
- Step 7: Improve domain boundaries (aggregates, invariants, value objects).

Each step should end with compiling code and an example commit message.

## Output Formatting
Copilot responses should include:
- A short explanation of the proposed incremental change.
- Updated or newly created files with complete code blocks.
- A suggested commit message.
- Notes about trade-offs or impacts when relevant.

