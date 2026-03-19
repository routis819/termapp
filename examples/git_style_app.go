package main

import (
	"fmt"
	"os"

	"github.com/routis819/termapp"
)

type GitStage struct {
	termapp.BaseStage
}

func (s *GitStage) Prompt() string { return "git> " }

func (s *GitStage) Commands() map[string]termapp.Command {
	return map[string]termapp.Command{
		"remote": {
			Description: "Manage set of tracked repositories",
			SubCommands: map[string]termapp.Command{
				"add": {
					Description: "Adds a remote",
					Handler: func(app *termapp.App, args []string) error {
						if len(args) < 2 {
							return fmt.Errorf("usage: remote add <name> <url>")
						}
						fmt.Printf("Added remote %s (%s)\n", args[0], args[1])
						return nil
					},
				},
				"remove": {
					Description: "Removes a remote",
					Completer: func(app *termapp.App, args []string) []string {
						return []string{"origin", "upstream", "backup"}
					},
					Handler: func(app *termapp.App, args []string) error {
						if len(args) < 1 {
							return fmt.Errorf("usage: remote remove <name>")
						}
						fmt.Printf("Removed remote %s\n", args[0])
						return nil
					},
				},
			},
		},
		"commit": {
			Description: "Record changes to the repository",
			Handler: func(app *termapp.App, args []string) error {
				for i, arg := range args {
					if arg == "-m" && i+1 < len(args) {
						fmt.Printf("Committed with message: %s\n", args[i+1])
						return nil
					}
				}
				return fmt.Errorf("usage: commit -m <message>")
			},
			Completer: func(app *termapp.App, args []string) []string {
				// Suggest common prefixes for commit messages if -m is present
				for i, arg := range args {
					if arg == "-m" && i == len(args)-2 {
						return []string{"feat: ", "fix: ", "docs: ", "chore: "}
					}
				}
				return []string{"-m"}
			},
		},
	}
}

func main() {
	app := termapp.NewApp(&GitStage{})
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
