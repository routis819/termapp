# Implementation Plan: Core Interaction Loop Tests

## Phase 1: Test Infrastructure and Basic Loop Tests
- [x] Task: Set up test infrastructure for mocking `liner` state. (b779543)
    - [x] Task: Write Tests: Create a test that attempts to run `App.Run` with a mocked input. (b779543)
    - [x] Task: Implement: Refactor `App` or use an interface if necessary to allow injecting a mocked `liner.State` or input stream. (b779543)
- [ ] Task: Test basic command dispatching in the loop.
    - [ ] Task: Write Tests: Verify that a command entered in the loop calls the expected handler.
    - [ ] Task: Implement: Ensure the test setup correctly feeds input and captures output.
- [ ] Task: Conductor - User Manual Verification 'Phase 1: Test Infrastructure and Basic Loop Tests' (Protocol in workflow.md)

## Phase 2: Lifecycle and Termination Tests
- [ ] Task: Verify lifecycle hooks are called during the `Run` loop.
    - [ ] Task: Write Tests: Check if `OnEnter` of the root stage is called.
    - [ ] Task: Implement: Update mock stage to record hook calls.
- [ ] Task: Test application termination.
    - [ ] Task: Write Tests: Verify `exit` command and `io.EOF` terminate the loop.
    - [ ] Task: Implement: Provide terminating input in tests.
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Lifecycle and Termination Tests' (Protocol in workflow.md)

## Phase 3: Coverage and Refinement
- [ ] Task: Verify and reach >80% coverage.
    - [ ] Task: Write Tests: Add tests for any uncovered branches in `Run`, `completer`, and `processCommand`.
    - [ ] Task: Implement: Run `go test -cover` and address gaps.
- [ ] Task: Conductor - User Manual Verification 'Phase 3: Coverage and Refinement' (Protocol in workflow.md)
