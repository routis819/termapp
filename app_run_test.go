package termapp

import (
	"io"
	"testing"

	"github.com/peterh/liner"
)

type mockInputter struct {
	inputs  []string
	history []string
	closed  bool
}

func (m *mockInputter) Prompt(prompt string) (string, error) {
	if len(m.inputs) == 0 {
		return "", io.EOF
	}
	input := m.inputs[0]
	m.inputs = m.inputs[1:]
	return input, nil
}

func (m *mockInputter) AppendHistory(item string) {
	m.history = append(m.history, item)
}

func (m *mockInputter) SetCompleter(f liner.Completer) {
}

func (m *mockInputter) Close() error {
	m.closed = true
	return nil
}

type flexibleMockStage struct {
	BaseStage
	prompt   string
	commands map[string]Command
}

func (f *flexibleMockStage) Prompt() string                { return f.prompt }
func (f *flexibleMockStage) Commands() map[string]Command { return f.commands }

func TestAppRun_Basic(t *testing.T) {
	mockS := &flexibleMockStage{prompt: "> "}
	inputter := &mockInputter{
		inputs: []string{"help", "exit"},
	}
	app := NewAppWithInputter(mockS, inputter)

	err := app.Run()
	if err != nil {
		t.Fatalf("App.Run failed: %v", err)
	}

	if !inputter.closed {
		t.Error("expected inputter to be closed")
	}
}

func TestAppRun_CommandDispatching(t *testing.T) {
	called := false
	mockS := &flexibleMockStage{
		prompt: "> ",
		commands: map[string]Command{
			"testcmd": {
				Description: "A test command",
				Handler: func(app *App, args []string) error {
					called = true
					return nil
				},
			},
		},
	}
	inputter := &mockInputter{
		inputs: []string{"testcmd", "exit"},
	}
	app := NewAppWithInputter(mockS, inputter)

	err := app.Run()
	if err != nil {
		t.Fatalf("App.Run failed: %v", err)
	}

	if !called {
		t.Error("expected command handler to be called")
	}
}
