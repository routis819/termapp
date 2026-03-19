package termapp

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

// Mock stage
type mockStage struct {
	BaseStage
	commands map[string]Command
}

func (m *mockStage) Prompt() string { return "> " }
func (m *mockStage) Commands() map[string]Command {
	if m.commands == nil {
		m.commands = map[string]Command{
			"testcmd": {
				Description: "A test command",
				Handler: func(app *App, args []string) error {
					return nil
				},
			},
		}
	}
	return m.commands
}

type conflictStage struct {
	BaseStage
}

func (s *conflictStage) Prompt() string { return "> " }
func (s *conflictStage) Commands() map[string]Command {
	return map[string]Command{
		"help": {
			Description: "Conflicting help",
			Handler:     func(app *App, args []string) error { return nil },
		},
	}
}

func captureStdout(f func()) string {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func captureStderr(f func()) string {
	r, w, _ := os.Pipe()
	oldStderr := os.Stderr
	os.Stderr = w

	f()

	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestProcessCommand_Help(t *testing.T) {
	app := NewApp(&mockStage{})

	output := captureStdout(func() {
		exit, err := app.processCommand("help")
		if exit || err != nil {
			t.Errorf("Expected continue without error, got exit=%v, err=%v", exit, err)
		}
	})

	if !strings.Contains(output, "Global commands:") {
		t.Errorf("Expected 'Global commands:' in output, got: %s", output)
	}
	if !strings.Contains(output, "testcmd") {
		t.Errorf("Expected 'testcmd' in output, got: %s", output)
	}
}

func TestProcessCommand_HelpWithArgs(t *testing.T) {
	app := NewApp(&mockStage{})

	output := captureStdout(func() {
		exit, err := app.processCommand("help extra arg")
		if exit || err != nil {
			t.Errorf("Expected continue without error, got exit=%v, err=%v", exit, err)
		}
	})

	if !strings.Contains(output, "Global commands:") {
		t.Errorf("Expected 'Global commands:' in output, got: %s", output)
	}
}

func TestProcessCommand_ExitAndQuit(t *testing.T) {
	app := NewApp(&mockStage{})

	output := captureStdout(func() {
		exit, err := app.processCommand("exit")
		if !exit || err != nil {
			t.Errorf("Expected exit=true, err=nil, got exit=%v, err=%v", exit, err)
		}
	})

	if !strings.Contains(output, "Exiting...") {
		t.Errorf("Expected 'Exiting...' in output, got: %s", output)
	}

	output = captureStdout(func() {
		exit, err := app.processCommand("quit extra args")
		if !exit || err != nil {
			t.Errorf("Expected exit=true, err=nil, got exit=%v, err=%v", exit, err)
		}
	})

	if !strings.Contains(output, "Exiting...") {
		t.Errorf("Expected 'Exiting...' in output, got: %s", output)
	}
}

func TestCompleter(t *testing.T) {
	app := NewApp(&mockStage{})

	// Test prefix matching global commands
	suggestions := app.Completer("he")
	if !reflect.DeepEqual(suggestions, []string{"help"}) {
		t.Errorf("Expected [help], got %v", suggestions)
	}

	suggestions = app.Completer("ex")
	if !reflect.DeepEqual(suggestions, []string{"exit"}) {
		t.Errorf("Expected [exit], got %v", suggestions)
	}

	// Test prefix matching stage commands
	suggestions = app.Completer("te")
	if !reflect.DeepEqual(suggestions, []string{"testcmd"}) {
		t.Errorf("Expected [testcmd], got %v", suggestions)
	}

	// Empty input should suggest all
	suggestions = app.Completer("")
	expected := []string{"help", "exit", "quit", "testcmd"}

	for _, exp := range expected {
		found := false
		for _, act := range suggestions {
			if exp == act {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected %s in suggestions, got %v", exp, suggestions)
		}
	}
}

func TestValidateCommands(t *testing.T) {
	app := NewApp(&mockStage{})

	stderr := captureStderr(func() {
		app.Push(&conflictStage{})
	})

	if !strings.Contains(stderr, "Warning: Stage command \"help\" conflicts") {
		t.Errorf("Expected conflict warning in stderr, got: %q", stderr)
	}
}

func TestSetRemoveGlobal(t *testing.T) {
	app := NewApp(&mockStage{})

	// Test SetGlobal
	app.SetGlobal("greet", Command{
		Description: "Greets user",
		Handler: func(app *App, args []string) error {
			fmt.Print("Hello!")
			return nil
		},
	})

	output := captureStdout(func() {
		app.processCommand("greet")
	})
	if output != "Hello!" {
		t.Errorf("Expected 'Hello!', got %q", output)
	}

	// Test override
	app.SetGlobal("help", Command{
		Description: "Custom help",
		Handler: func(app *App, args []string) error {
			fmt.Print("Custom Help Output")
			return nil
		},
	})
	output = captureStdout(func() {
		app.processCommand("help")
	})
	if output != "Custom Help Output" {
		t.Errorf("Expected 'Custom Help Output', got %q", output)
	}

	// Test RemoveGlobal
	app.RemoveGlobal("help")
	output = captureStdout(func() {
		app.processCommand("help")
	})
	if !strings.Contains(output, "Unknown command: help") {
		t.Errorf("Expected unknown command for help, got %q", output)
	}
}

func TestCompleter_CommandSpecific(t *testing.T) {
	app := NewApp(&mockStage{})

	// Add a command with a custom completer to the stage
	app.Current().Commands()["pick"] = Command{
		Description: "Pick a color",
		Handler:     func(app *App, args []string) error { return nil },
		Completer: func(app *App, args []string) []string {
			colors := []string{"red", "green", "blue"}
			var suggestions []string
			prefix := ""
			if len(args) > 0 {
				prefix = args[0]
			}
			for _, c := range colors {
				if strings.HasPrefix(c, prefix) {
					suggestions = append(suggestions, c)
				}
			}
			return suggestions
		},
	}

	// Test prefix matching the command-specific completer
	// Input: "pick r" -> args: ["r"], partial token: "r"
	// Expected suggestions: ["pick red"]
	// (Note: The current Completer returns full lines)
	suggestions := app.Completer("pick r")
	expected := []string{"pick red"}
	if !reflect.DeepEqual(suggestions, expected) {
		t.Errorf("Expected %v, got %v", expected, suggestions)
	}
}

func TestTokenize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "spaces only",
			input:    "    line ",
			expected: []string{"line"},
		},
		{
			name:     "multiple args",
			input:    "    line arg1     ",
			expected: []string{"line", "arg1"},
		},
		{
			name:     "double quotes",
			input:    `line "arg with spaces" arg2`,
			expected: []string{"line", "arg with spaces", "arg2"},
		},
		{
			name:     "single quotes",
			input:    `line 'arg with spaces' arg2`,
			expected: []string{"line", "arg with spaces", "arg2"},
		},
		{
			name:     "mixed quotes",
			input:    `line "single ' quote inside" 'double " quote inside'`,
			expected: []string{"line", "single ' quote inside", `double " quote inside`},
		},
		{
			name:     "escaped quotes",
			input:    `line \"escaped\"`,
			expected: []string{"line", `"escaped"`},
		},
		{
			name:     "empty quotes",
			input:    `line "" ''`,
			expected: []string{"line", "", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := tokenize(tt.input)
			if !reflect.DeepEqual(tokens, tt.expected) {
				t.Errorf("tokenize(%q) = %v, want %v", tt.input, tokens, tt.expected)
			}
		})
	}
}
