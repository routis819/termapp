# Product Definition: termapp

**termapp** is a Go-based framework for building stateful, interactive terminal applications. It manages context through a "Stage Stack" and provides automated, context-aware auto-completion.

## Core Vision
- **Context-Awareness**: The application knows exactly which "Stage" (screen) the user is in and adjusts its behavior (commands, auto-completion, prompt) accordingly.
- **Declarative Simplicity**: Developers define stages and commands using simple interfaces, and the framework handles the complex interaction loop and navigation.
- **Familiar Metaphors**: Uses the "Back Stack" (Push/Pop) concept from mobile/web development for terminal navigation.

## Key Goals
1.  **Stage Management**: Push and Pop stages to navigate between different application states.
2.  **Command-Handler Mapping**: Automatically route parsed input to stage-specific handlers.
3.  **Dynamic Completion**: Automatically update the `readline`/`liner` completer based on the active stage.
4.  **Lifecycle Management**: Provide hooks for entering and exiting stages (e.g., `OnEnter`, `OnExit`).
