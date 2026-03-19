// Package termapp provides a framework for building stateful, interactive terminal applications.
//
// It is designed to bridge the gap between low-level Readline wrappers and heavy TUI
// frameworks. It provides a structured, "stage-based" approach to build interactive shells.
//
// # Key Concepts
//
// Stage:
// A Stage represents a specific state or context of the application (e.g., Login, Dashboard).
// Each stage defines its own prompt, available commands, and lifecycle hooks.
//
// Back Stack:
// Navigation is handled via a stack of Stages. You can Push a new stage onto the stack
// to enter a new context, and Pop the current stage to return to the previous context,
// optionally passing a result back.
//
// Lifecycle Hooks:
// Stages can implement OnEnter, OnExit, OnDestroy, and OnResult hooks to manage setup,
// teardown, and data passing during navigation. The BaseStage struct can be embedded
// to provide empty default implementations for these hooks, reducing boilerplate.
//
// Global Commands:
// The framework automatically handles global commands like 'help', 'exit', and 'quit'
// across all stages, providing a consistent user experience. Auto-completion dynamically
// updates based on the active Stage and these global commands.
//
// Example:
//
//	type RootStage struct {
//	    termapp.BaseStage
//	}
//
//	func (s *RootStage) Prompt() string { return "root> " }
//	func (s *RootStage) Commands() map[string]termapp.Command {
//	    return map[string]termapp.Command{
//	        "hello": {
//	            Description: "Prints a greeting",
//	            Handler: func(app *termapp.App, args []string) error {
//	                fmt.Println("Hello, World!")
//	                return nil
//	            },
//	        },
//	    }
//	}
//
//	func main() {
//	    app := termapp.NewApp(&RootStage{})
//	    app.Run()
//	}
package termapp

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"

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
	// Completer provides dynamic completion candidates for the command's arguments.
	Completer func(app *App, args []string) []string
	// SubCommands for static hierarchical completion.
	SubCommands map[string]Command
}

// CommandInputter defines the operations for command-line input.
type CommandInputter interface {
	Prompt(prompt string) (string, error)
	AppendHistory(item string)
	SetCompleter(f liner.Completer)
	Close() error
}

// App orchestrates the lifecycle and the input loop.
type App struct {
	line    CommandInputter
	stack   []Stage
	Globals map[string]Command
}

// ErrExit is a special error that signals the application should terminate.
var ErrExit = fmt.Errorf("exit")

// NewApp creates a new termapp application with an initial stage.
func NewApp(root Stage) *App {
	l := liner.NewLiner()
	l.SetCtrlCAborts(true)
	l.SetTabCompletionStyle(liner.TabPrints)

	app := &App{
		line:  l,
		stack: []Stage{root},
	}
	app.initGlobals()
	return app
}

// NewAppWithInputter creates a new termapp application with a custom inputter.
func NewAppWithInputter(root Stage, inputter CommandInputter) *App {
	app := &App{
		line:  inputter,
		stack: []Stage{root},
	}
	app.initGlobals()
	return app
}

func (a *App) initGlobals() {
	a.Globals = map[string]Command{
		"help": {
			Description: "Show this help message",
			Handler: func(app *App, args []string) error {
				fmt.Println("Global commands:")
				for name, cmd := range app.Globals {
					fmt.Printf("  %s - %s\n", name, cmd.Description)
				}
				if curr := app.Current(); curr != nil && len(curr.Commands()) > 0 {
					fmt.Println("\nStage commands:")
					for name, cmd := range curr.Commands() {
						fmt.Printf("  %s - %s\n", name, cmd.Description)
					}
				}
				return nil
			},
		},
		"exit": {
			Description: "Exit the application",
			Handler: func(app *App, args []string) error {
				fmt.Println("Exiting...")
				return ErrExit
			},
		},
		"quit": {
			Description: "Exit the application",
			Handler: func(app *App, args []string) error {
				fmt.Println("Exiting...")
				return ErrExit
			},
		},
	}
}

// SetGlobal adds or overrides a global command.
func (a *App) SetGlobal(name string, cmd Command) {
	a.Globals[strings.ToLower(name)] = cmd
}

// RemoveGlobal removes a global command.
func (a *App) RemoveGlobal(name string) {
	delete(a.Globals, strings.ToLower(name))
}

func (a *App) validateCommands(s Stage) {
	if s == nil {
		return
	}
	for name := range s.Commands() {
		if _, ok := a.Globals[strings.ToLower(name)]; ok {
			fmt.Fprintf(os.Stderr, "Warning: Stage command %q conflicts with a global command. Global command takes precedence.\n", name)
		}
	}
}

// Push adds a new stage to the stack.
func (a *App) Push(s Stage) error {
	a.validateCommands(s)
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
		a.validateCommands(root)
		if err := root.OnEnter(a); err != nil {
			return fmt.Errorf("failed to enter root stage: %w", err)
		}
	}

	// Configure dynamic completion
	a.line.SetCompleter(a.Completer)

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

// Completer provides the default tab-completion logic for the application.
// It suggests global commands followed by commands from the active stage.
// It supports hierarchical completion via SubCommands and dynamic completion via Completer.
func (a *App) Completer(line string) []string {
	full, partial, state := tokenizeForCompletion(line)
	curr := a.Current()

	// Initial set of commands to check
	var candidates map[string]Command
	if curr != nil {
		candidates = curr.Commands()
	}

	// Traverse the command hierarchy as far as possible
	var currentCmd *Command
	var breadcrumbs []string
	var cmdLevel int

	for i, token := range full {
		found := false
		tokenLower := strings.ToLower(token)

		if cmdLevel == i {
			// Look in current candidates
			if cmd, ok := candidates[tokenLower]; ok {
				currentCmd = &cmd
				candidates = cmd.SubCommands
				breadcrumbs = append(breadcrumbs, token)
				cmdLevel++
				found = true
			} else if i == 0 {
				// Only check Globals for the first token
				if cmd, ok := a.Globals[tokenLower]; ok {
					currentCmd = &cmd
					candidates = cmd.SubCommands
					breadcrumbs = append(breadcrumbs, token)
					cmdLevel++
					found = true
				}
			}
		}

		if !found {
			// This token is not part of the command hierarchy, it's an argument.
			// We stop traversing subcommands, but we keep the token for reconstruction.
			break
		}
	}

	var suggestions []string
	prefix := strings.ToLower(partial)

	// Reconstruct the base line up to the partial token.
	// We MUST use the original tokens to preserve their formatting (quotes/spaces).
	// For simplicity, we use the tokens returned by the tokenizer, but we need to
	// be careful about quoting them if they contain spaces.
	var baseBuilder strings.Builder
	for _, token := range full {
		if strings.Contains(token, " ") {
			baseBuilder.WriteString("\"" + token + "\" ")
		} else {
			baseBuilder.WriteString(token + " ")
		}
	}
	baseLine := baseBuilder.String()

	// Suggest subcommands
	if len(full) == cmdLevel {
		for name := range candidates {
			if strings.HasPrefix(strings.ToLower(name), prefix) {
				suggestions = append(suggestions, baseLine+name)
			}
		}
	}

	// If we are at a specific command, check its dynamic completer
	if currentCmd != nil {
		// Arguments to the command are everything after cmdLevel
		args := full[cmdLevel:]
		args = append(args, partial)

		if currentCmd.Completer != nil {
			cmdSuggestions := currentCmd.Completer(a, args)
			for _, s := range cmdSuggestions {
				if strings.HasPrefix(strings.ToLower(s), prefix) {
					formatted := s
					if state == StateInDoubleQuote {
						formatted = "\"" + s + "\""
					} else if state == StateInSingleQuote {
						formatted = "'" + s + "'"
					} else if strings.Contains(s, " ") {
						// Auto-quote if it contains spaces and we are not in quotes
						formatted = "\"" + s + "\""
					}
					suggestions = append(suggestions, baseLine+formatted)
				}
			}
		}
	}

	// Also check Globals if we are at the root level
	if len(full) == 0 {
		for name := range a.Globals {
			if strings.HasPrefix(strings.ToLower(name), prefix) {
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

	curr := a.Current()
	var candidates map[string]Command
	if curr != nil {
		candidates = curr.Commands()
	}

	// Traverse the hierarchy to find the command to execute
	var cmd *Command
	var args []string

	for i, token := range tokens {
		tokenLower := strings.ToLower(token)
		found := false

		if i == 0 {
			// Check globals only at root
			if c, ok := a.Globals[tokenLower]; ok {
				cmd = &c
				candidates = c.SubCommands
				found = true
			}
		}

		if !found && candidates != nil {
			if c, ok := candidates[tokenLower]; ok {
				cmd = &c
				candidates = c.SubCommands
				found = true
			}
		}

		if found {
			// If we found a command, remaining tokens are potential arguments
			// unless we find a deeper subcommand.
			args = tokens[i+1:]
		} else {
			// Token not found in current candidates, it's an argument to the last found command
			break
		}
	}

	if cmd != nil {
		err := cmd.Handler(a, args)
		if err == ErrExit {
			return true, nil
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		return false, nil
	}

	fmt.Printf("Unknown command: %s. Type 'help' if available.\n", tokens[0])
	return false, nil
}

// CompletionState represents the parser state for command-line completion.
type CompletionState int

const (
	// StateNormal indicates the parser is between tokens or at a new token.
	StateNormal CompletionState = iota
	// StateInSingleQuote indicates the parser is inside a single-quoted string.
	StateInSingleQuote
	// StateInDoubleQuote indicates the parser is inside a double-quoted string.
	StateInDoubleQuote
)

func tokenizeForCompletion(line string) ([]string, string, CompletionState) {
	var full []string
	var current strings.Builder
	state := StateNormal
	escaped := false

	for _, r := range line {
		if escaped {
			current.WriteRune(r)
			escaped = false
			continue
		}

		if r == '\\' {
			escaped = true
			continue
		}

		switch state {
		case StateNormal:
			if r == '"' {
				state = StateInDoubleQuote
			} else if r == '\'' {
				state = StateInSingleQuote
			} else if unicode.IsSpace(r) {
				if current.Len() > 0 {
					full = append(full, current.String())
					current.Reset()
				}
			} else {
				current.WriteRune(r)
			}
		case StateInDoubleQuote:
			if r == '"' {
				full = append(full, current.String())
				current.Reset()
				state = StateNormal
			} else {
				current.WriteRune(r)
			}
		case StateInSingleQuote:
			if r == '\'' {
				full = append(full, current.String())
				current.Reset()
				state = StateNormal
			} else {
				current.WriteRune(r)
			}
		}
	}

	partial := current.String()
	return full, partial, state
}

func tokenize(line string) []string {
	var tokens []string
	var currentToken strings.Builder
	inSingleQuote := false
	inDoubleQuote := false
	escaped, inToken := false, false

	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}

	for _, r := range line {
		if escaped {
			currentToken.WriteRune(r)
			escaped = false
			inToken = true
			continue
		}

		if r == '\\' {
			escaped = true
			inToken = true
			continue
		}

		if r == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
			inToken = true
			continue
		}

		if r == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
			inToken = true
			continue
		}

		if (r == ' ' || r == '\t') && !inSingleQuote && !inDoubleQuote {
			if inToken {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
				inToken = false
			}
			continue
		}

		inToken = true
		currentToken.WriteRune(r)
	}

	if inToken {
		tokens = append(tokens, currentToken.String())
	}

	return tokens
}
