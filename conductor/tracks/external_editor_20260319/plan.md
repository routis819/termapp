# Implementation Plan: External Editor Stage (track-20260319-external-editor)

## Phase 1: Editor Interface and Core Logic
- [ ] Task: Define `Editor` interface and `SystemEditor` implementation.
    - [ ] Write Tests: Verify `SystemEditor` correctly identifies the editor from `$EDITOR` or fallbacks.
    - [ ] Implement: Create `editor.go` with `Editor` interface and `Launch` logic using `os/exec`.
- [ ] Task: Implement `ExternalEditorStage`.
    - [ ] Write Tests: Verify `ExternalEditorStage.OnEnter` triggers the editor launch.
    - [ ] Implement: Create `editor_stage.go` with lifecycle hooks managing terminal suspension and editor execution.
- [ ] Task: Conductor - User Manual Verification 'Phase 1: Editor Interface and Core Logic' (Protocol in workflow.md)

## Phase 2: Terminal and File Management
- [ ] Task: Integrate terminal suspension into `App`.
    - [ ] Write Tests: Verify `App` can pause and resume the `liner` inputter.
    - [ ] Implement: Update `App` and `CommandInputter` to support `Pause()` and `Resume()` if necessary, or manage closure/re-creation of `liner`.
- [ ] Task: Support managed temporary files.
    - [ ] Write Tests: Verify `ExternalEditorStage` creates and cleans up temporary files.
    - [ ] Implement: Use `os.CreateTemp` for managing the editing session and ensure cleanup in `OnDestroy`.
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Terminal and File Management' (Protocol in workflow.md)

## Phase 3: Data Passing and Validation
- [ ] Task: Implement `OnResult` data passing.
    - [ ] Write Tests: Verify the parent stage receives the edited content as a `string`.
    - [ ] Implement: Ensure `ExternalEditorStage.OnExit` reads the file and calls `app.Pop(content)`.
- [ ] Task: Add example to `simple_app.go`.
    - [ ] Implement: Update `examples/simple_app.go` with a command to launch the editor and display the result.
- [ ] Task: Conductor - User Manual Verification 'Phase 3: Data Passing and Validation' (Protocol in workflow.md)

## Phase 4: Final Verification
- [ ] Task: Comprehensive Integration Test.
    - [ ] Write Tests: Mock `os/exec` to simulate an editor session and verify the full flow from `Push` to `OnResult`.
    - [ ] Implement: Add `editor_test.go`.
- [ ] Task: Conductor - User Manual Verification 'Phase 4: Final Verification' (Protocol in workflow.md)
