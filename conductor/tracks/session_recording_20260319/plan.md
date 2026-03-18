# Implementation Plan: Session Recording and Framework Logging (track-20260319-session-recording)

## Phase 1: Core Logging & Output Infrastructure
- [ ] Task: Define JSON log entry structures and basic Logger.
    - [ ] Write Tests: Verify `Logger` can write and flush a JSONL entry to a file.
    - [ ] Implement: Create `logger.go` with JSON structs and `FileLogger` implementation.
- [ ] Task: Implement `StructuredWriter` for stage output orchestration.
    - [ ] Write Tests: Verify `StructuredWriter` captures output and attaches metadata (StageID, Timestamp).
    - [ ] Implement: Create `writer.go` with `StructuredWriter` implementation.
- [ ] Task: Integrate `Logger` and `StructuredWriter` into `App` struct.
    - [ ] Write Tests: Verify `App` can initialize with or without logging/structured output.
    - [ ] Implement: Update `App` and `NewApp` to include optional `Logger` and default `StructuredWriter` (stdout fallback).
- [ ] Task: Conductor - User Manual Verification 'Phase 1: Core Logging & Output Infrastructure' (Protocol in workflow.md)

## Phase 2: Framework Event Logging
- [ ] Task: Log stage transitions (Push/Pop).
    - [ ] Write Tests: Verify `Push` and `Pop` actions generate log entries.
    - [ ] Implement: Add logging calls to `App.Push` and `App.Pop`.
- [ ] Task: Log lifecycle hooks.
    - [ ] Write Tests: Verify `OnEnter`, `OnExit`, `OnDestroy`, `OnResult` calls generate logs.
    - [ ] Implement: Add logging calls to the lifecycle execution logic.
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Framework Event Logging' (Protocol in workflow.md)

## Phase 3: Session Recording (Inputs/Outputs)
- [ ] Task: Intercept and log user inputs.
    - [ ] Write Tests: Verify all inputs from `CommandInputter` are logged.
    - [ ] Implement: Update `App.Run` to log input tokens.
- [ ] Task: Log application outputs via `StructuredWriter`.
    - [ ] Write Tests: Verify that output written to `StructuredWriter` generates log entries.
    - [ ] Implement: Configure `StructuredWriter` to pipe output to the active `Logger` when recording.
- [ ] Task: Conductor - User Manual Verification 'Phase 3: Session Recording (Inputs/Outputs)' (Protocol in workflow.md)

## Phase 4: Global Commands and CLI Flags
- [ ] Task: Implement `record` global command.
    - [ ] Write Tests: Verify `record start <file>` initiates recording and `record stop` terminates it.
    - [ ] Implement: Add `record` command handler to global command interception.
- [ ] Task: Support CLI flag recording.
    - [ ] Write Tests: Verify recording starts immediately if specified in configuration.
    - [ ] Implement: Update `NewApp` to accept recording configuration from CLI flags or options.
- [ ] Task: Conductor - User Manual Verification 'Phase 4: Global Commands and CLI Flags' (Protocol in workflow.md)

## Phase 5: Final Validation and Documentation
- [ ] Task: Comprehensive Integration Test.
    - [ ] Write Tests: Run a full session with multiple stages and verify the integrity of the generated JSONL file.
    - [ ] Implement: Add `session_recording_test.go`.
- [ ] Task: Update example.
    - [ ] Implement: Update `examples/simple_app.go` to support a `--record` flag for demonstration.
- [ ] Task: Conductor - User Manual Verification 'Phase 5: Final Validation and Documentation' (Protocol in workflow.md)
