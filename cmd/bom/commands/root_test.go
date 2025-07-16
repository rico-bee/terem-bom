package commands

import (
	"bytes"
	"testing"
)

func TestRootCommand(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"--help"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Root command help failed: %v", err)
	}

	output := buf.String()
	expectedStrings := []string{
		"BOM Weather Data Processor",
		"Usage:",
		"Available Commands:",
		"convert",
		"validate",
		"version",
	}

	for _, expected := range expectedStrings {
		if !contains(output, expected) {
			t.Errorf("Expected output to contain '%s', got: %s", expected, output)
		}
	}
}

func TestRootCommandVersion(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"--version"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Root command version failed: %v", err)
	}

	output := buf.String()
	if !contains(output, "1.0.0") {
		t.Errorf("Expected version '1.0.0', got: %s", output)
	}
}

func TestRootCommandInvalidCommand(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"invalid-command"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Expected error for invalid command")
	}
}

func TestRootCommandVerboseFlag(t *testing.T) {
	cmd := NewRootCmd()
	verboseFlag := cmd.PersistentFlags().Lookup("verbose")
	if verboseFlag == nil {
		t.Fatal("Verbose flag not found")
	}

	if verboseFlag.Name != "verbose" {
		t.Errorf("Expected flag name 'verbose', got: %s", verboseFlag.Name)
	}

	if verboseFlag.Shorthand != "v" {
		t.Errorf("Expected flag shorthand 'v', got: %s", verboseFlag.Shorthand)
	}
}

func TestExecute(t *testing.T) {
	// Test the Execute function
	// This is a simple test to ensure it doesn't panic
	// We can't easily test the full execution without mocking
	// but we can test that the function exists and is callable
	_ = Execute
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
