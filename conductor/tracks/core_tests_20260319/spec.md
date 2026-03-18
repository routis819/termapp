# Specification: Core Interaction Loop Tests

This track focuses on increasing the test coverage and reliability of the core `App.Run` interaction loop.

## Objectives
- Verify that `App.Run` correctly initializes the active stage.
- Verify that `App.Run` handles user input and dispatches it correctly.
- Verify that `App.Run` terminates correctly on exit commands or EOF.
- Achieve >80% coverage for the `Run` and `processCommand` methods.

## Technical Design
- Use `io.Pipe` or a custom `io.ReadWriter` to mock terminal input/output for `liner`.
- Create mock stages to verify lifecycle hooks (`OnEnter`, `OnExit`) are called during the loop.
- Test edge cases like empty input, unknown commands, and global command interception.

## Success Criteria
- [ ] Unit tests cover the main `Run` loop.
- [ ] Test coverage for `app.go` exceeds 80%.
- [ ] No regressions introduced in existing functionality.
