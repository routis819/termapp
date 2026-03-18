# Tech Stack: termapp

## Programming Language
- **Go (Golang)**: Chosen for its performance, simplicity, and excellent standard library for CLI tools.

## Core Dependencies
- **github.com/peterh/liner**: A pure Go line-editing library. It is chosen for its simplicity and ease of implementing dynamic auto-completion.
- **regexp**: Used for POSIX-compliant tokenization of command arguments.

## Design Patterns
- **Stage Pattern**: Each app state is an implementation of a `Stage` interface.
- **Stack Pattern**: Navigation is managed via a Last-In-First-Out (LIFO) stack of Stages.
- **Lifecycle Events**: Each stage implements standard lifecycle hooks:
    - `OnEnter`: Called when the stage becomes the active top of the stack.
    - `OnExit`: Called when the stage is no longer the top (either via Push or Pop).
    - `OnDestroy`: Called when a stage is permanently removed from the stack (Pop).
    - `OnResult`: Called when a previous stage (pushed onto the current one) is Popped and returns data.
- **Command Dispatcher**: Input is tokenized and dispatched to stage-specific command maps.
- **Global Command Interception**: The framework intercepts specific global commands (e.g., `help`, `exit`, `quit`) before dispatching stage-specific commands, ensuring a consistent user experience.
