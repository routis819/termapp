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
	// Not needed for basic run test
}

func (m *mockInputter) Close() error {
	m.closed = true
	return nil
}

func TestAppRun_Basic(t *testing.T) {
	mockS := &mockStage{}
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

	// Verify help was called (we capture output in app_test.go's captureStdout but app.Run prints to stdout)
	// We might need to refactor app to allow injecting an output writer too if we want to verify stdout.
}
