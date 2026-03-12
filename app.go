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
}

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

// Pop removes the current stage and returns to the previous one.
func (a *App) Pop() error {
	if len(a.stack) <= 1 {
		return fmt.Errorf("already at the root stage")
	}

	current := a.Current()
	if err := current.OnExit(a); err != nil {
		return err
	}

	a.stack = a.stack[:len(a.stack)-1]

	next := a.Current()
	if next != nil {
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
	a.line.SetCompleter(func(line string) []string {
		curr := a.Current()
		if curr == nil {
			return nil
		}

		var suggestions []string
		for name := range curr.Commands() {
			if strings.HasPrefix(name, strings.ToLower(line)) {
				suggestions = append(suggestions, name)
			}
		}
		return suggestions
	})

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
		tokens := tokenize(input)
		if len(tokens) == 0 {
			continue
		}

		cmdName := tokens[0]
		args := tokens[1:]

		if cmd, ok := curr.Commands()[cmdName]; ok {
			if err := cmd.Handler(a, args); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}
		} else {
			fmt.Printf("Unknown command: %s. Type 'help' if available.\n", cmdName)
		}
	}
	return nil
}

func tokenize(line string) []string {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}
	re := regexp.MustCompile(`\s+`)
	return re.Split(line, -1)
}
