# Specification: Global Commands (Help/Exit)

This track implements a unified way to handle 'help' and 'exit' commands across all application stages.

## Objectives
- Standardize 'help' output for all stages.
- Provide a global 'exit' command to terminate the application from any stage.
- Centralize command dispatching to handle these global keywords before stage-specific logic.

## Functional Requirements
1. **Help Command**:
   - Typing `help` should display a list of all available commands for the current stage.
   - The output should include the command name and a brief description (if available).
2. **Exit Command**:
   - Typing `exit` (or `quit`) should terminate the application loop regardless of the current stage stack depth.

## Technical Design
- Update the `App` or `Manager` to intercept `help` and `exit` inputs.
- The `Stage` interface already has `Commands() map[string]Command`. We will iterate over this map for the `help` command.
- We might need a `Command` struct if it's currently just a function, to include a description.

## Success Criteria
- [ ] `help` command works in the example app and displays correct commands.
- [ ] `exit` command terminates the app from any stage.
- [ ] Unit tests verify that global commands take precedence or are handled consistently.
