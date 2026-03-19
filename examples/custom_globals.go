package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/peterh/liner"
	"github.com/routis819/termapp"
)

type MyStage struct {
	termapp.BaseStage
}

func (s *MyStage) Prompt() string { return "custom> " }
func (s *MyStage) Commands() map[string]termapp.Command {
	return map[string]termapp.Command{
		"hello": {
			Description: "Prints hello",
			Handler: func(app *termapp.App, args []string) error {
				fmt.Println("Hello from stage!")
				return nil
			},
		},
	}
}

func main() {
	app := termapp.NewApp(&MyStage{})

	// Override exit to ask for confirmation
	app.SetGlobal("exit", termapp.Command{
		Description: "Exit with confirmation",
		Handler: func(a *termapp.App, args []string) error {
			l := liner.NewLiner()
			defer l.Close()
			ans, _ := l.Prompt("Are you sure you want to exit? (y/n): ")
			if strings.ToLower(ans) == "y" {
				fmt.Println("Goodbye!")
				return termapp.ErrExit
			}
			return nil
		},
	})

	// Add a new global command 'whoami'
	app.SetGlobal("whoami", termapp.Command{
		Description: "Prints user info",
		Handler: func(a *termapp.App, args []string) error {
			fmt.Printf("User: %s\n", os.Getenv("USER"))
			return nil
		},
	})

	// Remove 'quit' to force use of 'exit'
	app.RemoveGlobal("quit")

	fmt.Println("Custom Globals Example")
	fmt.Println("Try 'whoami', 'help', 'quit' (should fail), and 'exit'.")
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
