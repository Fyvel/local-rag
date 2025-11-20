# Copilot Code Cleanup Instructions (Go)

## Goal
Audit the Go codebase to identify and eliminate dead code, simplify redundant logic, enforce early-exit patterns, and rename modules or functions that do not follow Go conventions.  
All cleanup must preserve behavior, maintain compilation, and avoid introducing unnecessary abstractions.

## Scope of Cleanup

### 1. Remove unused code
Copilot should:
- Detect functions, types, structs, interfaces, methods, and modules that are never referenced.
- Remove unused constants or variables.
- Remove obsolete packages left from previous refactors.
- Confirm that removing an item will not break external interfaces or public APIs unless the user explicitly approves.

Constraints:
- If the unused item might be part of a public API, flag it rather than removing it silently.
- Perform removals incrementally, one commit per logical deletion.

### 2. Simplify redundant functions
Copilot should:
- Find functions whose only behavior is to call another function directly without adding value.
- Replace unnecessary wrappers with direct calls.
- Delete trivial pass-through abstractions unless they provide interface boundaries or dependency inversion.
- Collapse redundant helper functions into a single, clear implementation.

Trade-off rule:
- If a simplified function served as a seam for testing or an abstraction boundary, propose simplification but explain the impact.

### 3. Prefer early exits
Copilot should:
- Identify functions using deep indentation or nested `if` blocks.
- Rewrite using early returns to improve readability and reduce cyclomatic complexity.
- Keep early exits consistent with error-handling idioms:
  - `if err != nil { return … }`
- Ensure no semantic change.

### 4. Rename modules and functions to match Go conventions
Copilot should:
- Rename items not following Go naming rules:
  - Package names: all lowercase, no underscores, short, domain-specific.
  - Function names: CamelCase, exported only when needed.
  - No stutter (e.g., `user.UserService` → `user.Service`).
- Update all import paths and references.
- Ensure renames are incremental and compilable.

Examples:
- `helpers_util` → `util`
- `Calculate_total_price` → `CalculateTotalPrice`
- `ordersDomain` → `order`
- `DoStuff` → rename to a meaningful action

### 5. Preserve behavior and maintain build health
Copilot must:
- Keep the build green after every cleanup step.
- Avoid breaking imports.
- Remove or update tests when functions are removed or renamed.
  
If a change is risky, Copilot must:
- Propose the change,
- Explain the concern,
- Wait for explicit approval if needed.

## Output Format for Each Cleanup Step
Each suggestion from Copilot must include:

1. **Short explanation**  
   What was detected, why it is unused or redundant, or which naming guideline applies.

2. **Updated or new files (full code blocks)**  
   Always provide fully rewritten file content for modifications.

3. **List of deletions or renames**  
   Include old path and new path for renamed items.

4. **Suggested commit message**  
   Clear, concise, referencing the specific cleanup.

5. **Trade-off notes**  
   Only if simplification or removal could affect extensibility, testing, or architectural boundaries.
