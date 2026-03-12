# Implementation Plan: Enhanced Lifecycle (track-001-ext)

## Phase 1: Core Lifecycle Enhancement
- [ ] Update `Stage` interface in `app.go`.
- [ ] Add `BaseStage` struct to provide default empty implementations (to avoid boilerplate).
- [ ] Refactor `App.Pop` to accept `result interface{}`.
- [ ] Ensure `OnDestroy` is called only when a stage is truly removed.
- [ ] Ensure `OnResult` is called on the parent stage before it re-enters.

## Phase 2: Example Update
- [ ] Update `examples/simple_app.go` to use `BaseStage`.
- [ ] Create a "Selection" stage that returns a value to the "Main" stage via `OnResult`.

## Phase 3: Validation
- [ ] Verify that data is correctly passed between stages.
- [ ] Verify that `OnDestroy` is not called during a `Push`.
