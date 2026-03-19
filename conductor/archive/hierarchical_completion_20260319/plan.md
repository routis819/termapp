# Implementation Plan: Hierarchical Autocompletion with Quoted String Support

## Phase 1: Dynamic Candidate Generation (Foundation)
1.  **Enhance `Command` struct in `app.go`**: [DONE]
    -   Add `Completer func(app *App, args []string) []string`.
2.  **Update `App.completer`**: [DONE]
    -   Initial implementation to check if the first token is a command and, if so, call its `Completer`.
    -   This provides immediate value for simple commands with dynamic arguments (e.g., `pick [color]`).
3.  **Add Test Case**: [DONE]
    -   Verify that `pick <TAB>` returns dynamic suggestions from the `Completer` function.

## Phase 2: Completion-Aware Tokenizer
1.  **Implement `tokenizeForCompletion(line string) ([]string, string, State)`**: [DONE]
    -   `[]string`: Completed tokens.
    -   `string`: The current partial token at the end of the line.
    -   `State`: Contextual state (e.g., `InDoubleQuotes`).
2.  **Integrate Quoted String Support**: [DONE]
    -   Ensure the tokenizer tracks quote state to allow completion *inside* quotes.
3.  **Add Unit Tests**: [DONE]
    -   Test `cmd "arg with sp` -> `["cmd"], "arg with sp", InDoubleQuotes`.

## Phase 3: Hierarchical Completer Implementation
1.  **Refactor `App.completer` and `processCommand`**: [DONE]
    -   Add `SubCommands map[string]Command` to `Command` struct.
    -   Update `App.completer` to traverse the command tree.
    -   Update `processCommand` to handle nested command execution.
2.  **Implement Quoted Suggestion Formatting**: [DONE]
    -   If `InDoubleQuotes`, ensure suggestions are correctly formatted (e.g., closing quotes).

## Phase 4: Verification & Examples
1.  **New Test Suite**: `completion_test.go` [DONE]
    -   Exhaustive tests for all completion scenarios.
2.  **New Example**: `examples/git_style_app.go` [DONE]
    -   Demonstrate nested subcommands and dynamic completion for arguments.
3.  **Documentation Update**: [DONE]
    -   Update `README.md` and `conductor/product.md`.

## Verification
-   `go test ./...`
-   Manual verification with `examples/git_style_app.go`.

## Phase: Review Fixes
- [x] Task: Perform code review
