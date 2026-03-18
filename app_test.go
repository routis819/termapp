package termapp

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"testing"
)

// Mock stage
type mockStage struct {
	BaseStage
}

func (m *mockStage) Prompt() string { return "> " }
func (m *mockStage) Commands() map[string]Command {
	return map[string]Command{
		"testcmd": {
			Description: "A test command",
			Handler: func(app *App, args []string) error {
				return nil
			},
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

func TestProcessCommand_Help(t *testing.T) {
	app := NewApp(&mockStage{})

	output := captureStdout(func() {
		exit, err := app.processCommand("help")
		if exit || err != nil {
			t.Errorf("Expected continue without error, got exit=%v, err=%v", exit, err)
		}
	})

	if !bytes.Contains([]byte(output), []byte("Global commands:")) {
		t.Errorf("Expected 'Global commands:' in output, got: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("testcmd")) {
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

	if !bytes.Contains([]byte(output), []byte("Global commands:")) {
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

	if !bytes.Contains([]byte(output), []byte("Exiting...")) {
		t.Errorf("Expected 'Exiting...' in output, got: %s", output)
	}

	output = captureStdout(func() {
		exit, err := app.processCommand("quit extra args")
		if !exit || err != nil {
			t.Errorf("Expected exit=true, err=nil, got exit=%v, err=%v", exit, err)
		}
	})

	if !bytes.Contains([]byte(output), []byte("Exiting...")) {
		t.Errorf("Expected 'Exiting...' in output, got: %s", output)
	}
}

func TestCompleter(t *testing.T) {
	app := NewApp(&mockStage{})

	// Test prefix matching global commands
	suggestions := app.completer("he")
	if !reflect.DeepEqual(suggestions, []string{"help"}) {
		t.Errorf("Expected [help], got %v", suggestions)
	}

	suggestions = app.completer("ex")
	if !reflect.DeepEqual(suggestions, []string{"exit"}) {
		t.Errorf("Expected [exit], got %v", suggestions)
	}

	// Test prefix matching stage commands
	suggestions = app.completer("te")
	if !reflect.DeepEqual(suggestions, []string{"testcmd"}) {
		t.Errorf("Expected [testcmd], got %v", suggestions)
	}

	// Empty input should suggest all
	suggestions = app.completer("")
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

func TestTokenizerSpace(t *testing.T) {
	line := "    line "
	tokens := tokenize(line)

	if len(tokens) != 1 {
		t.FailNow()
	}

	var arr0 [1]string
	copy(arr0[:], tokens)
	if [1]string{"line"} != arr0 {
		t.Fail()
	}

	line = "    line arg1     "
	tokens = tokenize(line)

	if len(tokens) != 2 {
		t.Fail()
	}

	var arr1 [2]string
	copy(arr1[:], tokens)
	if [2]string{"line", "arg1"} != arr1 {
		t.Fail()
	}
}
