# Specification: Enhanced Lifecycle (track-001-ext)

This extension track focuses on adding `OnDestroy` and `OnResult` lifecycle hooks to the `Stage` interface and the `App` manager.

## Objectives
- Introduce `OnDestroy` for resource cleanup when a stage is popped from the stack.
- Introduce `OnResult` to allow parent stages to receive data from child stages when they are popped.
- Refactor `App.Pop` to accept an optional result value.

## Interface Changes

### Stage Interface
```go
type Stage interface {
    Prompt() string
    Commands() map[string]Command
    OnEnter(app *App) error
    OnExit(app *App) error
    OnDestroy(app *App) error
    OnResult(app *App, result interface{}) error
}
```

## Behavior Changes

### Pop(result interface{})
1. Get the current top stage `S1`.
2. Get the stage below it `S2`.
3. Call `S1.OnExit(app)`.
4. Call `S1.OnDestroy(app)`.
5. Remove `S1` from the stack.
6. If `S2` exists:
   - Call `S2.OnResult(app, result)`.
   - Call `S2.OnEnter(app)`.

## Success Criteria
- [ ] `app.go` updated with new methods and `Pop` signature.
- [ ] `examples/simple_app.go` demonstrates `OnResult` data passing.
- [ ] `OnDestroy` logs are visible during cleanup.
