package termapp

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/peterh/liner"
)

// Stage defines the behavior for a specific application state.
type Stage interface {
	Prompt() string
	Commands() map[string]Command

	OnEnter(app *App) error
	OnExit(app *App) error
	OnDestroy(app *App) error
	OnResult(app *App, result any) error
}

// BaseStage provides default empty implementations for lifecycle hooks.
// Embed this in your stage struct to avoid boilerplate.
type BaseStage struct{}

func (s *BaseStage) OnEnter(app *App) error              { return nil }
func (s *BaseStage) OnExit(app *App) error               { return nil }
func (s *BaseStage) OnDestroy(app *App) error            { return nil }
func (s *BaseStage) OnResult(app *App, result any) error { return nil }

// Command maps a user input to a function.
type Command struct {
	Description string
	Handler     func(app *App, args []string) error
}

// App orchestrates the lifecycle and the liner loop.
type App struct {
	line  *liner.State
	stack []Stage
}

// NewApp creates a new termapp application with an initial stage.
func NewApp(root Stage) *App {
	l := liner.NewLiner()
	l.SetCtrlCAborts(true)
	l.SetTabCompletionStyle(liner.TabPrints)

	return &App{
		line:  l,
		stack: []Stage{root},
	}
}

// Push adds a new stage to the stack.
func (a *App) Push(s Stage) error {
	if current := a.Current(); current != nil {
		if err := current.OnExit(a); err != nil {
			return err
		}
	}

	a.stack = append(a.stack, s)

	if err := s.OnEnter(a); err != nil {
		// If OnEnter fails, remove the stage from stack to maintain consistency
		a.stack = a.stack[:len(a.stack)-1]
		return err
	}
	return nil
}

// Pop removes the current stage and returns to the previous one, passing a result.
func (a *App) Pop(result any) error {
	if len(a.stack) <= 1 {
		return fmt.Errorf("already at the root stage")
	}

	current := a.Current()
	if err := current.OnExit(a); err != nil {
		return err
	}
	if err := current.OnDestroy(a); err != nil {
		return err
	}

	a.stack = a.stack[:len(a.stack)-1]

	next := a.Current()
	if next != nil {
		if err := next.OnResult(a, result); err != nil {
			return err
		}
		if err := next.OnEnter(a); err != nil {
			return err
		}
	}
	return nil
}

// Current returns the active stage.
func (a *App) Current() Stage {
	if len(a.stack) == 0 {
		return nil
	}
	return a.stack[len(a.stack)-1]
}

// Run starts the main interaction loop.
func (a *App) Run() error {
	defer a.line.Close()

	// Initial OnEnter for the root stage
	if root := a.Current(); root != nil {
		if err := root.OnEnter(a); err != nil {
			return fmt.Errorf("failed to enter root stage: %w", err)
		}
	}

	// Configure dynamic completion
	a.line.SetCompleter(a.completer)

	for {
		curr := a.Current()
		if curr == nil {
			break
		}

		input, err := a.line.Prompt(curr.Prompt())
		if err != nil {
			if err == liner.ErrPromptAborted || err == io.EOF {
				fmt.Println("\nExiting...")
				return nil
			}
			return err
		}

		if strings.TrimSpace(input) == "" {
			continue
		}

		a.line.AppendHistory(input)
		exit, err := a.processCommand(input)
		if exit {
			return err
		}
	}
	return nil
}

func (a *App) completer(line string) []string {
	curr := a.Current()
	var suggestions []string

	globalCmds := []string{"help", "exit", "quit"}
	for _, cmd := range globalCmds {
		if strings.HasPrefix(cmd, strings.ToLower(line)) {
			suggestions = append(suggestions, cmd)
		}
	}

	if curr != nil {
		for name := range curr.Commands() {
			if strings.HasPrefix(name, strings.ToLower(line)) {
				suggestions = append(suggestions, name)
			}
		}
	}
	return suggestions
}

func (a *App) processCommand(input string) (bool, error) {
	tokens := tokenize(input)
	if len(tokens) == 0 {
		return false, nil
	}

	cmdName := strings.ToLower(tokens[0])
	args := tokens[1:]

	curr := a.Current()

	// Global commands
	switch cmdName {
	case "help":
		fmt.Println("Global commands:")
		fmt.Println("  help - Show this help message")
		fmt.Println("  exit, quit - Exit the application")
		if curr != nil && len(curr.Commands()) > 0 {
			fmt.Println("\nStage commands:")
			for name, cmd := range curr.Commands() {
				fmt.Printf("  %s - %s\n", name, cmd.Description)
			}
		}
		return false, nil
	case "exit", "quit":
		fmt.Println("Exiting...")
		return true, nil
	}

	if curr != nil {
		if cmd, ok := curr.Commands()[cmdName]; ok {
			if err := cmd.Handler(a, args); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}
			return false, nil
		}
	}

	fmt.Printf("Unknown command: %s. Type 'help' if available.\n", cmdName)
	return false, nil
}

func tokenize(line string) []string {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}
	re := regexp.MustCompile(`\s+`)
	return re.Split(line, -1)
}
