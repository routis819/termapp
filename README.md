# termapp

**termapp** is a Go-based framework designed for building stateful, interactive terminal applications. It draws inspiration from mobile and web navigation patterns—such as Android's `Activity`/`Fragment` stacks or web page routing—to manage the complex state of a command-line interface.

## 🎯 Vision

Traditional CLI libraries in Go often fall into two extremes:
1. **Low-level Readline Wrappers**: Provide basic input handling but require developers to manually manage state transitions, command parsing, and global variables.
2. **Heavy TUI Frameworks**: Provide full-screen graphics (TUI) but introduce high complexity and move away from the standard "shell-like" stream-of-text experience.

**termapp** bridges this gap. It provides a structured, "stage-based" approach to build interactive shells where:
- Each screen or context is a self-contained **Stage**.
- Navigation is handled via a **Back Stack** (Push/Pop).
- **Auto-completion** is context-aware, automatically updating based on the active Stage.
- **Command Dispatching** is simplified through declarative mapping.

---

## 🏗 Design Architecture

### 1. The "Stage" Concept
A **Stage** represents a specific state of the application (e.g., `LoginStage`, `DashboardStage`, `SettingsStage`). Each stage defines its own:
- **Prompt**: What the user sees (e.g., `admin@dashboard> `).
- **Commands**: A localized set of available actions and their handlers.
- **Auto-completion**: A list of suggestions relevant only to that stage.
- **Lifecycle Events**: Hooks to manage state transitions:
    - `OnEnter`: Called when the stage becomes the active top of the stack.
    - `OnExit`: Called when the stage is no longer the top.
    - `OnDestroy`: Called when a stage is permanently removed from the stack.
    - `OnResult`: Called when a returning stage passes data back to this one.

### 2. Navigation Stack
The framework maintains a stack of Stages.
- **Push**: Enter a new sub-context (e.g., from Dashboard to Edit User).
- **Pop**: Return to the previous context (Activity-like "Back" behavior).
- **Home**: Clear the stack and return to the root Stage.

### 3. Context-Aware Completion
Unlike static completion, **termapp** dynamically reconfigures the `liner` completer whenever a Stage transition occurs. This ensures that users are only suggested commands that are valid in the current context.

---

## 🛠 Core API

```go
// Stage defines the behavior for a specific application state.
type Stage interface {
    Prompt() string
    Commands() map[string]Command
    OnEnter(app *App) error
    OnExit(app *App) error
    OnDestroy(app *App) error
    OnResult(app *App, result interface{}) error
}

// Command maps a user input to a function.
type Command struct {
    Description string
    Handler     func(app *App, args []string) error
}

// App orchestrates the lifecycle and the interaction loop.
type App struct {
    stack []Stage
    // ... internals
}
```

## 🚀 Current Status
- [x] Basic `liner` loop integration.
- [x] POSIX-compliant tokenization for command arguments.
- [x] Centralized Stage/App architecture.
- [x] Dynamic auto-completion provider.
- [ ] Built-in "Help" and "Exit" command management.

## 📝 Usage Example (Mockup)

```go
func main() {
    app := termapp.NewApp(DefaultStage{})
    
    // The App automatically handles the loop and delegates
    // input to the top stage in the stack.
    app.Run()
}
```

---

## ⚖️ License
MIT
