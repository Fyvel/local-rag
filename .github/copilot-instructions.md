# Copilot Architecture Review Instructions (Go + DDD Compliance)

## Goal
Review and validate the current architecture to ensure full compliance with the project's Domain-Driven Design (DDD) structure:

- `api/`  
- `application/`  
- `domain/`  
- `infra/`  
- `di/`

Identify misuse, misplaced dependencies, architectural leaks, or package violations, and propose corrective refactor steps.  
Finally, ensure the root `README.md` accurately documents the resulting architecture and code layout.

## Review Scope

Copilot should actively examine:
1. **Package structure**
   - Files placed in the wrong layer.
   - Packages not following Go naming standards.
   - Missing or redundant subpackages.

2. **Imports and dependency flows**
   - `domain` must not import `application`, `api`, `infra`, or `di`.
   - `application` must not import `api` or `di`.
   - `api` must not import `infra` directly (except when strictly necessary for low-level integration).
   - `di` may import everything but must not contain business logic.

3. **Misplaced logic**
   - Business rules in handlers → should be in `domain`.
   - Orchestration logic in domain → should be in `application`.
   - Infrastructure concerns in domain → move to `infra`.

4. **Circular dependencies**
   - Identify cycles and propose breakage strategies (e.g., extracting interfaces, restructuring packages).

5. **Inconsistent abstractions**
   - Unnecessary interfaces.
   - Repositories defined in `infra` instead of `domain`.
   - Infra leaking concrete types into higher layers.

6. **Tech debt from migration**
   - Stale helper functions.
   - Duplicate models.
   - Old folder remnants.
   - Forgotten TODOs referencing previous architecture.

## Correction Rules

Copilot must propose iterative fixes that:
- Are small and self-contained.
- Compile successfully.
- Preserve behavior.
- Improve structural correctness.

Valid types of correction steps:
- Moving structs, functions, or packages.
- Introducing minimal refactor interfaces to break cycles.
- Replacing infra types with domain abstractions.
- Consolidating domain entities and value objects.
- Standardizing naming (e.g., `order/`, not `ordersDomain/`).

## README Update Requirements

After architectural cleanup:
- Generate the updated `README.md`.
- It must describe the final DDD structure in plain language.
- It must include:
  - Folder overview (`api`, `application`, `domain`, `infra`, `di`)
  - High-level responsibilities for each layer
  - Explanation of dependency flow
  - Notes on Go package conventions used by the project
  - Steps for running, testing, and extending the project (if detectable)
- Rewrite sections that reflect old structure or concepts.

The README must not include obsolete or transitional information.

## Output Format

For each Copilot suggestion:
1. Short explanation of detected issue.
2. Corrected file changes (with full updated content).
3. Files moved or removed.
4. Updated README contents (when applicable).
5. Suggested commit message.
6. Brief trade-off notes (if relevant).


