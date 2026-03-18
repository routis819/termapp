# Specification: Session Recording and Framework Logging (track-20260319-session-recording)

## Overview
This feature introduces a comprehensive session recording and framework logging system for `termapp`. It allows users to capture full sessions (inputs and outputs) and framework-level events (stage transitions, errors) into a structured JSONL (JSON Lines) file on disk. Recording can be initiated via global commands or application-level CLI flags.

## Functional Requirements
- **Session Recording**:
    - Capture all user inputs received via the `CommandInputter`.
    - Capture all application outputs.
    - Store each interaction as a discrete JSON object on a new line (JSONL format).
- **Framework Logging**:
    - Log framework-level events, including:
        - Stage transitions (`Push`, `Pop`).
        - Lifecycle hooks (`OnEnter`, `OnExit`, `OnDestroy`, `OnResult`).
        - Framework errors and warnings.
    - Log entries should follow the same JSONL format for consistency.
- **Control Mechanisms**:
    - **Global Commands**: Implement `record start [filename]` and `record stop` as global commands available in all stages.
    - **CLI Flag Support**: Provide a mechanism for `termapp`-based applications to enable recording on startup (e.g., via a `--record` flag or a configuration option in `NewApp`).
- **Persistence**:
    - Automatically create and append to the specified output file on disk.
    - Ensure data is flushed to disk regularly to prevent loss on crash.

## Non-Functional Requirements
- **Performance**: The recording/logging mechanism should have minimal impact on the application's responsiveness.
- **Robustness**: Error-handling for file I/O should be robust, ensuring that a failure to write logs doesn't crash the application.
- **Extensibility**: The JSON schema for logs should be extensible to support future event types.

## Acceptance Criteria
- A session can be started using a global command `record start <filename>`.
- All inputs and outputs are recorded to the specified file in JSONL format.
- Framework events (like stage pushing/popping) are logged to the same file.
- The session recording can be stopped using `record stop`.
- `termapp` applications can enable recording via a CLI flag or programmatic configuration.
- The output file contains valid JSON objects, one per line.

## Out of Scope
- Automated playback of recorded sessions (this may be a future track).
- Real-time streaming of logs to remote servers.
- Complex log rotation policies.
