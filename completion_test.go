package termapp

import (
	"reflect"
	"testing"
)

func TestTokenizeForCompletion(t *testing.T) {
	tests := []struct {
		name          string
		line          string
		expectedFull  []string
		expectedPart  string
		expectedState CompletionState
	}{
		{
			"Empty line",
			"",
			nil,
			"",
			StateNormal,
		},
		{
			"Single command, no space",
			"help",
			nil,
			"help",
			StateNormal,
		},
		{
			"Single command with space",
			"help ",
			[]string{"help"},
			"",
			StateNormal,
		},
		{
			"Command and partial argument",
			"pick r",
			[]string{"pick"},
			"r",
			StateNormal,
		},
		{
			"Multiple arguments",
			"cmd arg1 arg2",
			[]string{"cmd", "arg1"},
			"arg2",
			StateNormal,
		},
		{
			"Starting double quotes",
			`cmd "`,
			[]string{"cmd"},
			"",
			StateInDoubleQuote,
		},
		{
			"Inside double quotes",
			`cmd "arg with space`,
			[]string{"cmd"},
			"arg with space",
			StateInDoubleQuote,
		},
		{
			"Closing double quotes",
			`cmd "quoted arg"`,
			[]string{"cmd", "quoted arg"},
			"",
			StateNormal,
		},
		{
			"Starting single quotes",
			"cmd '",
			[]string{"cmd"},
			"",
			StateInSingleQuote,
		},
		{
			"Inside single quotes",
			"cmd 'arg with space",
			[]string{"cmd"},
			"arg with space",
			StateInSingleQuote,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			full, part, state := tokenizeForCompletion(tt.line)
			if !reflect.DeepEqual(full, tt.expectedFull) {
				t.Errorf("Full tokens mismatch: expected %v, got %v", tt.expectedFull, full)
			}
			if part != tt.expectedPart {
				t.Errorf("Partial token mismatch: expected %q, got %q", tt.expectedPart, part)
			}
			if state != tt.expectedState {
				t.Errorf("State mismatch: expected %v, got %v", tt.expectedState, state)
			}
		})
	}
}
