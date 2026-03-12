# Implementation Plan: Core Refactoring & Lifecycle (track-001)

## Phase 1: Core Refactoring
- [ ] Update `Stage` interface in `app.go`.
- [ ] Refactor `App.Push` and `App.Pop` to handle lifecycle hooks.
- [ ] Update `NewApp` to be cleaner (lazy initialization or explicit start).

## Phase 2: Example Update
- [ ] Update `examples/simple_app.go` to implement the new `OnEnter`/`OnExit` methods.
- [ ] Demonstrate state initialization in `OnEnter`.

## Phase 3: Validation
- [ ] Run the example and verify the logs for lifecycle transitions.
- [ ] Fix any issues with `liner` state management during transitions.
