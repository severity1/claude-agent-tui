// Package accessibility provides utilities for accessibility compliance.
package accessibility

import "os"

// ReduceMotion returns true if the user prefers reduced motion.
// Checks the REDUCE_MOTION environment variable for accessibility compliance.
// When true, components should skip or minimize animations.
func ReduceMotion() bool {
	return os.Getenv("REDUCE_MOTION") != ""
}
