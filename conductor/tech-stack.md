# Tech Stack: termapp

## Programming Language
- **Go (Golang)**: Chosen for its performance, simplicity, and excellent standard library for CLI tools.

## Core Dependencies
- **github.com/peterh/liner**: A pure Go line-editing library. It is chosen for its simplicity and ease of implementing dynamic auto-completion.
- **regexp**: Used for POSIX-compliant tokenization of command arguments.

## Design Patterns
- **Stage Pattern**: Each app state is an implementation of a `Stage` interface.
- **Stack Pattern**: Navigation is managed via a Last-In-First-Out (LIFO) stack of Stages.
- **Command Dispatcher**: Input is tokenized and dispatched to stage-specific command maps.
