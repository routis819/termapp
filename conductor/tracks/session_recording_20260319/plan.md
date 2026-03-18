# Implementation Plan: Session Recording and Framework Logging (track-20260319-session-recording)

## Phase 1: Core Logging Infrastructure
- [ ] Task: Define JSON log entry structures and basic Logger.
    - [ ] Write Tests: Verify `Logger` can write and flush a JSONL entry to a file.
    - [ ] Implement: Create `logger.go` with JSON structs and `FileLogger` implementation.
- [ ] Task: Integrate `Logger` into `App` struct.
    - [ ] Write Tests: Verify `App` can initialize with or without a logger.
    - [ ] Implement: Update `App` and `NewApp` to include optional `Logger` support.
- [ ] Task: Conductor - User Manual Verification 'Phase 1: Core Logging Infrastructure' (Protocol in workflow.md)

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
- [ ] Task: Capture and log application outputs.
    - [ ] Write Tests: Verify that output from handlers (e.g., printed text) is captured and logged.
    - [ ] Implement: Introduce an output wrapper or interceptor for capturing stage responses.
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
