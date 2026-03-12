package main

import (
	"fmt"
	"github.com/routis819/termapp"
)

// RootStage: The main entry point of the app.
type RootStage struct{}

func (s *RootStage) Prompt() string {
	return "root> "
}

func (s *RootStage) OnEnter(app *termapp.App) error {
	fmt.Println("[RootStage] Entering... (Initializing state)")
	return nil
}

func (s *RootStage) OnExit(app *termapp.App) error {
	fmt.Println("[RootStage] Exiting...")
	return nil
}

func (s *RootStage) Commands() map[string]termapp.Command {
	return map[string]termapp.Command{
		"enter": {
			Description: "Enter the sub-stage",
			Handler: func(app *termapp.App, args []string) error {
				fmt.Println("Navigating to SubStage...")
				return app.Push(&SubStage{})
			},
		},
		"hello": {
			Description: "Say hello",
			Handler: func(app *termapp.App, args []string) error {
				fmt.Println("Hello from RootStage!")
				return nil
			},
		},
	}
}

// SubStage: A nested context with different commands.
type SubStage struct{}

func (s *SubStage) Prompt() string {
	return "sub-stage# "
}

func (s *SubStage) OnEnter(app *termapp.App) error {
	fmt.Println("[SubStage] Welcome! (Stage-specific setup)")
	return nil
}

func (s *SubStage) OnExit(app *termapp.App) error {
	fmt.Println("[SubStage] Goodbye! (Cleaning up)")
	return nil
}

func (s *SubStage) Commands() map[string]termapp.Command {
	return map[string]termapp.Command{
		"back": {
			Description: "Go back to the previous stage",
			Handler: func(app *termapp.App, args []string) error {
				fmt.Println("Returning to previous stage...")
				return app.Pop()
			},
		},
		"secret": {
			Description: "A secret command only in SubStage",
			Handler: func(app *termapp.App, args []string) error {
				fmt.Println("You found the secret! Arguments passed:", args)
				return nil
			},
		},
	}
}

func main() {
	app := termapp.NewApp(&RootStage{})

	fmt.Println("--- termapp Lifecycle & Navigation Example ---")
	fmt.Println("Commands are now mapped directly to handlers with access to the App instance.")
	fmt.Println("----------------------------------------------")

	if err := app.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
