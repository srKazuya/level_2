package builtins_test

import (
	"github.com/stretchr/testify/assert"
	"minishell/internal/utils"
	"minishell/internal/commands/builtins"
	"testing"
)

func TestEcho(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "basic echo",
			args:     []string{"echo", "hello", "world"},
			expected: "hello world\n",
		},
		{
			name:     "no newline (-n)",
			args:     []string{"echo", "-n", "hello"},
			expected: "hello",
		},
		{
			name:     "enable escapes (-e)",
			args:     []string{"echo", "-e", "line1\\nline2"},
			expected: "line1\nline2\n",
		},
		{
			name:     "tabs and backslash (-e)",
			args:     []string{"echo", "-e", "a\\tb\\\\c"},
			expected: "a\tb\\c\n",
		},
		{
			name:     "empty args",
			args:     []string{"echo"},
			expected: "\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output := utils.CaptureOutput(func() {
				builtins.Echo(tc.args)
			})
			assert.Equal(t, tc.expected, output)
		})
	}
}

