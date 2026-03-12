# Implementation Plan: Global Commands (Help/Exit)

## Phase 1: Global Command Interception
- [ ] Task: Implement `Help` handler in `App.Run`.
    - [ ] Task: Write Tests: Create `app_test.go` (if not present/adequate) and verify `help` command triggers expected output.
    - [ ] Task: Implement: Update `App.Run` loop to intercept `help` before stage-specific commands.
- [ ] Task: Implement `Exit` handler in `App.Run`.
    - [ ] Task: Write Tests: Verify `exit` command terminates the loop.
    - [ ] Task: Implement: Add logic to handle `exit` and `quit`.
- [ ] Task: Conductor - User Manual Verification 'Phase 1: Global Command Interception' (Protocol in workflow.md)

## Phase 2: Refinement and Auto-completion
- [ ] Task: Add Global Commands to Auto-completer.
    - [ ] Task: Write Tests: Verify `help` and `exit` appear in tab completion.
    - [ ] Task: Implement: Update the completer function in `App.Run`.
- [ ] Task: Ensure robust tokenization for global commands.
    - [ ] Task: Write Tests: Test `help <arg>` and `exit <arg>` behavior.
    - [ ] Task: Implement: Handle extra arguments gracefully (e.g., ignore them).
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Refinement and Auto-completion' (Protocol in workflow.md)
