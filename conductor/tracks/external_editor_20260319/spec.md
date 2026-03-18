# Specification: External Editor Stage (track-20260319-external-editor)

## Overview
This feature introduces an `ExternalEditorStage` to `termapp`, allowing applications to launch a system text editor (e.g., `vim`, `nano`) to capture multi-line input or edit files. The framework will handle the complexities of pausing the `liner` inputter, launching the editor process, managing temporary files, and returning the result to the parent stage.

## Functional Requirements
- **Editor Selection**:
    - Use the `$EDITOR` environment variable to determine the editor.
    - Fallback to common editors (e.g., `vi` or `nano`) if `$EDITOR` is not set.
- **Terminal Management**:
    - Automatically suspend the `liner` inputter (raw mode) before launching the editor.
    - Restore the terminal state and resume the inputter after the editor process exits.
- **File & Content Management**:
    - Support creating and automatically cleaning up managed temporary files.
    - Provide a mechanism to "seed" the editor with initial content.
    - Support editing existing files by path (optional/future-proof).
- **Data Passing**:
    - Return the full content of the edited file to the parent stage via the `OnResult` lifecycle hook.
    - Ensure the return type is flexible (e.g., `string` or `[]byte`).

## Non-Functional Requirements
- **Reliability**: Ensure the terminal is always restored to a usable state, even if the editor crashes or is interrupted.
- **User Experience**: The transition between the `termapp` prompt and the external editor should be seamless.
- **Platform Compatibility**: The implementation should work across common Linux and Unix-like environments.

## Acceptance Criteria
- Pushing an `ExternalEditorStage` successfully launches the system editor.
- The user can edit content in the editor and save/exit.
- Upon editor exit, the terminal returns to the `termapp` prompt.
- The parent stage receives the edited content in its `OnResult` hook.
- Temporary files are properly cleaned up after use.

## Out of Scope
- Built-in multi-line text editor within `termapp` (this feature relies on external tools).
- Complex terminal multiplexing or split-pane support.
