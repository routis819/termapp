package termapp

import (
	"errors"
	"io"
	"testing"

	"github.com/peterh/liner"
)

type mockInputter struct {
	inputs      []string
	history     []string
	closed      bool
	promptErr   error
	onCompleter func(f liner.Completer)
}

func (m *mockInputter) Prompt(prompt string) (string, error) {
	if m.promptErr != nil {
		return "", m.promptErr
	}
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
	if m.onCompleter != nil {
		m.onCompleter(f)
	}
}

func (m *mockInputter) Close() error {
	m.closed = true
	return nil
}

type flexibleMockStage struct {
	BaseStage
	prompt     string
	commands   map[string]Command
	onEnter    func(app *App) error
	onExit     func(app *App) error
	onResult   func(app *App, result any) error
	onDestroy  func(app *App) error
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
func (f *flexibleMockStage) OnResult(app *App, result any) error {
	if f.onResult != nil {
		return f.onResult(app, result)
	}
	return nil
}
func (f *flexibleMockStage) OnDestroy(app *App) error {
	if f.onDestroy != nil {
		return f.onDestroy(app)
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

func TestAppRun_Termination(t *testing.T) {
	t.Run("EOF", func(t *testing.T) {
		mockS := &flexibleMockStage{prompt: "> "}
		inputter := &mockInputter{
			inputs: []string{}, // Immediately EOF
		}
		app := NewAppWithInputter(mockS, inputter)
		err := app.Run()
		if err != nil {
			t.Errorf("expected no error on EOF, got %v", err)
		}
	})

	t.Run("ExitCommand", func(t *testing.T) {
		mockS := &flexibleMockStage{prompt: "> "}
		inputter := &mockInputter{
			inputs: []string{"exit"},
		}
		app := NewAppWithInputter(mockS, inputter)
		err := app.Run()
		if err != nil {
			t.Errorf("expected no error on exit, got %v", err)
		}
	})

	t.Run("Aborted", func(t *testing.T) {
		mockS := &flexibleMockStage{prompt: "> "}
		inputter := &mockInputter{
			promptErr: liner.ErrPromptAborted,
		}
		app := NewAppWithInputter(mockS, inputter)
		err := app.Run()
		if err != nil {
			t.Errorf("expected no error on aborted prompt, got %v", err)
		}
	})
}

func TestAppRun_EdgeCases(t *testing.T) {
	t.Run("EmptyInput", func(t *testing.T) {
		mockS := &flexibleMockStage{prompt: "> "}
		inputter := &mockInputter{
			inputs: []string{"", "  ", "\t", "exit"},
		}
		app := NewAppWithInputter(mockS, inputter)
		err := app.Run()
		if err != nil {
			t.Fatalf("App.Run failed: %v", err)
		}
		// If we didn't hang, it handled empty inputs correctly
	})

	t.Run("UnknownCommand", func(t *testing.T) {
		mockS := &flexibleMockStage{prompt: "> "}
		inputter := &mockInputter{
			inputs: []string{"unknown", "exit"},
		}
		app := NewAppWithInputter(mockS, inputter)
		err := app.Run()
		if err != nil {
			t.Fatalf("App.Run failed: %v", err)
		}
	})

	t.Run("PromptError", func(t *testing.T) {
		mockS := &flexibleMockStage{prompt: "> "}
		expectedErr := errors.New("prompt error")
		inputter := &mockInputter{
			promptErr: expectedErr,
		}
		app := NewAppWithInputter(mockS, inputter)
		err := app.Run()
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("RootOnEnterError", func(t *testing.T) {
		expectedErr := errors.New("enter error")
		mockS := &flexibleMockStage{
			prompt: "> ",
			onEnter: func(app *App) error {
				return expectedErr
			},
		}
		inputter := &mockInputter{}
		app := NewAppWithInputter(mockS, inputter)
		err := app.Run()
		if err == nil || !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})
}

func TestAppNavigation(t *testing.T) {
	root := &flexibleMockStage{prompt: "root> "}
	app := NewAppWithInputter(root, &mockInputter{})

	t.Run("PushSuccess", func(t *testing.T) {
		child := &flexibleMockStage{prompt: "child> "}
		err := app.Push(child)
		if err != nil {
			t.Fatalf("Push failed: %v", err)
		}
		if app.Current() != child {
			t.Error("expected current stage to be child")
		}
	})

	t.Run("PopSuccess", func(t *testing.T) {
		resultReceived := ""
		root.onResult = func(app *App, result any) error {
			resultReceived = result.(string)
			return nil
		}
		err := app.Pop("result-data")
		if err != nil {
			t.Fatalf("Pop failed: %v", err)
		}
		if app.Current() != root {
			t.Error("expected current stage to be root")
		}
		if resultReceived != "result-data" {
			t.Errorf("expected result-data, got %s", resultReceived)
		}
	})

	t.Run("PopAtRootError", func(t *testing.T) {
		err := app.Pop("result")
		if err == nil {
			t.Error("expected error when popping at root")
		}
	})

	t.Run("PushError", func(t *testing.T) {
		expectedErr := errors.New("push enter error")
		child := &flexibleMockStage{
			prompt: "child> ",
			onEnter: func(app *App) error {
				return expectedErr
			},
		}
		err := app.Push(child)
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
		if app.Current() != root {
			t.Error("expected current stage to remain root after failed push")
		}
	})

	t.Run("PopExitError", func(t *testing.T) {
		expectedErr := errors.New("exit error")
		child := &flexibleMockStage{
			prompt: "child> ",
			onExit: func(app *App) error {
				return expectedErr
			},
		}
		app.Push(child)
		err := app.Pop("result")
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("PopDestroyError", func(t *testing.T) {
		// Reset app
		app = NewAppWithInputter(root, &mockInputter{})
		expectedErr := errors.New("destroy error")
		child := &flexibleMockStage{
			prompt: "child> ",
			onDestroy: func(app *App) error {
				return expectedErr
			},
		}
		app.Push(child)
		err := app.Pop("result")
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("CurrentEmpty", func(t *testing.T) {
		emptyApp := &App{stack: []Stage{}}
		if emptyApp.Current() != nil {
			t.Error("expected nil for empty stack")
		}
	})
}
