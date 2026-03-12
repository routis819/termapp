package main

import (
	"fmt"
	"github.com/routis819/termapp"
	"strings"
)

// RootStage embeds BaseStage to inherit default lifecycle methods.
type RootStage struct {
	termapp.BaseStage
	lastResult string
}

func (s *RootStage) Prompt() string {
	prompt := "root> "
	if s.lastResult != "" {
		prompt = fmt.Sprintf("root(last:%s)> ", s.lastResult)
	}
	return prompt
}

// OnResult is called when a child stage pops and returns data.
func (s *RootStage) OnResult(app *termapp.App, result interface{}) error {
	if res, ok := result.(string); ok {
		s.lastResult = res
		fmt.Printf("[RootStage] Received result from child: %s\n", res)
	}
	return nil
}

func (s *RootStage) Commands() map[string]termapp.Command {
	return map[string]termapp.Command{
		"select": {
			Description: "Go to selection stage",
			Handler: func(app *termapp.App, args []string) error {
				return app.Push(&SelectionStage{})
			},
		},
		"exit": {
			Description: "Exit the app",
			Handler: func(app *termapp.App, args []string) error {
				fmt.Println("Goodbye!")
				return nil // In a real app, you might use a flag to stop the loop
			},
		},
	}
}

// SelectionStage allows user to pick a value.
type SelectionStage struct {
	termapp.BaseStage
}

func (s *SelectionStage) Prompt() string {
	return "choose-color [red/green/blue]> "
}

func (s *SelectionStage) OnDestroy(app *termapp.App) error {
	fmt.Println("[SelectionStage] Cleaning up resources before destruction...")
	return nil
}

func (s *SelectionStage) Commands() map[string]termapp.Command {
	return map[string]termapp.Command{
		"pick": {
			Description: "Pick a color and return",
			Handler: func(app *termapp.App, args []string) error {
				if len(args) == 0 {
					return fmt.Errorf("please provide a color")
				}
				color := strings.ToLower(args[0])
				fmt.Printf("[SelectionStage] Picking %s and returning to root...\n", color)
				return app.Pop(color) // Pass the result back to RootStage
			},
		},
		"cancel": {
			Description: "Cancel and return with no result",
			Handler: func(app *termapp.App, args []string) error {
				return app.Pop(nil)
			},
		},
	}
}

func main() {
	app := termapp.NewApp(&RootStage{})

	fmt.Println("--- termapp Enhanced Lifecycle Example ---")
	fmt.Println("1. Type 'select' to enter SelectionStage.")
	fmt.Println("2. In SelectionStage, type 'pick [color]' (e.g., 'pick red').")
	fmt.Println("3. Watch RootStage receive the result and update its prompt!")
	fmt.Println("------------------------------------------")

	if err := app.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
