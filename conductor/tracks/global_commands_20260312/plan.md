# Implementation Plan: Global Commands (Help/Exit)

## Phase 1: Global Command Interception [checkpoint: 3831fe0]
- [x] Task: Implement `Help` handler in `App.Run`. (12f79c2)
    - [x] Task: Write Tests: Create `app_test.go` (if not present/adequate) and verify `help` command triggers expected output. (12f79c2)
    - [x] Task: Implement: Update `App.Run` loop to intercept `help` before stage-specific commands. (12f79c2)
- [x] Task: Implement `Exit` handler in `App.Run`. (12f79c2)
    - [x] Task: Write Tests: Verify `exit` command terminates the loop. (12f79c2)
    - [x] Task: Implement: Add logic to handle `exit` and `quit`. (12f79c2)
- [x] Task: Conductor - User Manual Verification 'Phase 1: Global Command Interception' (3831fe0)

## Phase 2: Refinement and Auto-completion
- [x] Task: Add Global Commands to Auto-completer.
    - [x] Task: Write Tests: Verify `help` and `exit` appear in tab completion.
    - [x] Task: Implement: Update the completer function in `App.Run`.
- [x] Task: Ensure robust tokenization for global commands.
    - [x] Task: Write Tests: Test `help <arg>` and `exit <arg>` behavior.
    - [x] Task: Implement: Handle extra arguments gracefully (e.g., ignore them).
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Refinement and Auto-completion' (Protocol in workflow.md)
