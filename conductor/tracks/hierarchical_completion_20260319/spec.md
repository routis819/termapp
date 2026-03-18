# Specification: Hierarchical Autocompletion with Quoted String Support

## Objective
Enable `termapp` to provide context-aware, hierarchical autocompletion for subcommands and arguments, correctly handling quoted strings and whitespace.

## Background
Currently, `termapp` only completes the first word (the command) against stage-specific and global commands. To support more complex CLI interactions, the framework needs to:
1.  Understand the current command hierarchy (e.g., `git push ` should suggest remotes).
2.  Provide suggestions even inside quoted strings (e.g., `git commit -m "feat: `).
3.  Detect when a user has finished one token (with a space) and is ready for the next.

## Technical Requirements

### 1. Dynamic Candidate Generation
The framework must provide a mechanism for stages and commands to generate completion candidates dynamically (e.g., from a list of files, remote branches, or current state).
-   `Stage` or `Command` can provide a `CandidateFunc(app *App, args []string) []string`.
-   Candidates should be filtered by prefix and correctly formatted for `liner`.

### 2. Completion-Aware Tokenizer (`CompletionTokenizer`)
The current `tokenize` function in `app.go` is designed for command *execution*. It trims spaces and returns a final list of tokens. For *completion*, we need a more granular parser that:
-   Returns the state of the line: `InToken`, `BetweenTokens`, `InSingleQuote`, `InDoubleQuote`.
-   Preserves trailing whitespace information (does not use `strings.TrimSpace`).
-   Identifies the "partial token" currently being typed.

### 3. Enhanced `Command` Structure
The `Command` struct should be optionally recursive or support a dedicated completer.
```go
type Command struct {
    Description string
    Handler     func(app *App, args []string) error
    // New: SubCommands for static hierarchical completion
    SubCommands map[string]Command
    // New: DynamicCompleter for programmatic argument completion
    DynamicCompleter func(app *App, args []string) []string
}
```

### 3. Updated `App.completer` Logic
The `App.completer` (the function passed to `liner`) must:
1.  Use `CompletionTokenizer` to parse the line.
2.  Navigate the `Stage.Commands()` and `Command.SubCommands` maps based on completed tokens.
3.  If a `DynamicCompleter` is found for the current context, invoke it.
4.  Handle prefix matching for the "partial token" at the end of the line.

### 4. Quote Handling in Suggestions
If a user is completing inside a quote, the suggestions returned by `liner` should either include the quotes or be compatible with how `liner` replaces the current word. (Liner's `Completer` replaces the *entire* line up to the cursor, so we must return the full prefix + suggestion).

## Success Criteria
-   `git <TAB>` suggests subcommands.
-   `git pu<TAB>` completes to `git push`.
-   `git push <TAB>` calls a dynamic completer (if defined) to suggest remotes.
-   `cmd "arg with sp<TAB>` completes to `cmd "arg with space"`.
-   All existing tests pass.
-   New tests cover nested completion and quoted strings.
