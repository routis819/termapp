# Implementation Plan: Robust and Customizable Global Commands and Completion

## Phase 1: Data Structure & App Initialization
1.  **Modify `App` struct**:
    -   Add `Globals map[string]Command`.
2.  **Update `NewApp`**:
    -   Initialize `Globals` with default handlers for `help`, `exit`, and `quit`.
    -   Ensure these defaults use the same behavior as currently implemented.
3.  **Refactor `processCommand`**:
    -   Replace the hardcoded `switch` statement with a lookup in `a.Globals`.

## Phase 2: Refactoring Completer
1.  **Refactor `completer` method**:
    -   Instead of hardcoded strings, iterate over `a.Globals` to generate suggestions.
2.  **Expose Default Completer (Optional but Recommended)**:
    -   Provide a public method `DefaultCompleter(line string) []string` that can be called by custom completers.

## Phase 3: Validation & Conflict Detection
1.  **Implement `validateCommands(stage Stage)`**:
    -   Check for name collisions between `stage.Commands()` and `a.Globals`.
2.  **Integrate into Lifecycle**:
    -   Call `validateCommands` during `Push` and initial `Run`.
    -   Decide on behavior (warning or error) if a conflict is found.

## Phase 4: Customization API
1.  **Add `SetGlobal(name string, cmd Command)`**:
    -   Allow users to add or override global commands.
2.  **Add `RemoveGlobal(name string)`**:
    -   Allow users to remove global commands.

## Phase 5: Verification & Examples
1.  **Update `app_test.go`**:
    -   Verify that default globals still work.
    -   Verify that overriding a global works.
    -   Verify that removing a global works.
    -   Verify that validation detects conflicts.
2.  **New Example**: `examples/custom_globals.go`
    -   Showcase how to customize the "shell" experience.

## Verification
-   `go test ./...`
-   Manual verification of the new example.
