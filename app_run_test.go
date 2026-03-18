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
	onEnter  func(app *App) error
	onExit   func(app *App) error
}

func (f *flexibleMockStage) Prompt() string                { return f.prompt }
func (f *flexibleMockStage) Commands() map[string]Command { return f.commands }
func (f *flexibleMockStage) OnEnter(app *App) error {
	if f.onEnter != nil {
		return f.onEnter(app)
	}
	return nil
}
func (f *flexibleMockStage) OnExit(app *App) error {
	if f.onExit != nil {
		return f.onExit(app)
	}
	return nil
}

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

func TestAppRun_LifecycleHooks(t *testing.T) {
	onEnterCalled := false
	mockS := &flexibleMockStage{
		prompt: "> ",
		onEnter: func(app *App) error {
			onEnterCalled = true
			return nil
		},
	}
	inputter := &mockInputter{
		inputs: []string{"exit"},
	}
	app := NewAppWithInputter(mockS, inputter)

	err := app.Run()
	if err != nil {
		t.Fatalf("App.Run failed: %v", err)
	}

	if !onEnterCalled {
		t.Error("expected OnEnter to be called")
	}
}
