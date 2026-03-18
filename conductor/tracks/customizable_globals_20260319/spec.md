# Specification: Robust and Customizable Global Commands and Completion

## Objective
Provide a robust and flexible system for managing "Global Commands"
(like `help`, `exit`) and their corresponding autocompletion, allowing
developers to override, remove, or extend them.

## Current Limitations
1.  **Hardcoded Globals**: Commands like `help`, `exit`, and `quit`
    are hardcoded in `App.completer` and `App.processCommand`.
2.  **No Override**: Developers cannot change the behavior or name of
    these global commands.
3.  **Fragile Completion**: The completion logic is tightly coupled to
    these hardcoded strings, making it difficult to extend without
    modifying the core framework.
4.  **No Validation**: There is no check to see if a `Stage` defines a
    command that conflicts with a global one.

## Technical Requirements

### 1. Registry for Global Commands
-   The `App` struct should hold a map or registry of global commands.
-   These globals should be initialized with sensible defaults but be
    modifiable before `Run()`.

### 2. Global Command Definition
A global command should be defined similarly to stage-specific
commands, with a description and a handler.

### 3. Composable Completer
-   The framework should provide a way to obtain the "default" list of
    suggestions.
-   Developers should be able to wrap or replace the default completer
    logic while still having access to the framework's completion
    utility functions.

### 4. Validation & Conflict Detection
-   When `Run()` is called (or when a stage is pushed), the framework
    should ideally check for name collisions between stage-specific
    and global commands.
-   Global commands should take precedence by default, but this should
    be configurable.

## Success Criteria
-   The default `help`, `exit`, and `quit` behavior is preserved but
    no longer hardcoded.
-   A developer can remove the `quit` command entirely.
-   A developer can add a new global `status` command that works in
    all stages.
-   Autocompletion automatically reflects the current set of
    registered global commands.
-   A warning or error is issued if a stage command conflicts with a
    global command.
