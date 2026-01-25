package accessibility_test

import (
	"os"
	"testing"

	"github.com/severity1/claude-agent-tui/internal/accessibility"
)

func TestReduceMotion_EnvSet(t *testing.T) {
	// Set the environment variable
	os.Setenv("REDUCE_MOTION", "1")
	defer os.Unsetenv("REDUCE_MOTION")

	if !accessibility.ReduceMotion() {
		t.Error("ReduceMotion() = false, want true when REDUCE_MOTION is set")
	}
}

func TestReduceMotion_EnvUnset(t *testing.T) {
	// Ensure the environment variable is not set
	os.Unsetenv("REDUCE_MOTION")

	if accessibility.ReduceMotion() {
		t.Error("ReduceMotion() = true, want false when REDUCE_MOTION is not set")
	}
}

func TestReduceMotion_EmptyStringIsFalsy(t *testing.T) {
	// Set to empty string - should be falsy since os.Getenv returns ""
	os.Setenv("REDUCE_MOTION", "")
	defer os.Unsetenv("REDUCE_MOTION")

	// Empty string means "" != "" is false, so ReduceMotion returns false
	if accessibility.ReduceMotion() {
		t.Error("ReduceMotion() = true, want false when REDUCE_MOTION is empty string")
	}
}

func TestReduceMotion_AnyValueIsTruthy(t *testing.T) {
	testCases := []string{"true", "false", "0", "yes", "no", "anything"}

	for _, val := range testCases {
		os.Setenv("REDUCE_MOTION", val)

		if !accessibility.ReduceMotion() {
			t.Errorf("ReduceMotion() = false with REDUCE_MOTION=%q, want true", val)
		}
	}
	os.Unsetenv("REDUCE_MOTION")
}
