# Specification: Core Refactoring & Lifecycle (track-001)

This track focuses on formalizing the `Stage` lifecycle and improving the `App` manager to handle stage transitions more robustly.

## Objectives
- Introduce `OnEnter` and `OnExit` hooks to the `Stage` interface.
- Ensure `App.Push` and `App.Pop` trigger these hooks correctly.
- Refactor `App` to provide a way for stages to access the `App` instance without manual boilerplate.

## Interface Changes

### Stage Interface
```go
type Stage interface {
    Prompt() string
    Commands() map[string]Command
    OnEnter(app *App) error
    OnExit(app *App) error
}
```

## Behavior
1. **Push**:
   - Call `Current().OnExit(app)` for the old stage (if it exists).
   - Add new stage to the stack.
   - Call `NewStage().OnEnter(app)` for the new stage.
2. **Pop**:
   - Call `Current().OnExit(app)` for the stage being removed.
   - Remove stage from the stack.
   - Call `Current().OnEnter(app)` for the returning stage.
3. **Run**:
   - The loop starts and maintains the stack.

## Success Criteria
- [ ] `app.go` refactored with new interface.
- [ ] Lifecycle hooks are called during navigation.
- [ ] Example `simple_app.go` updated and running correctly.
