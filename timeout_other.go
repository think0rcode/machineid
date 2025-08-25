//go:build !windows

package machineid

import "time"

// Default timeout for non-Windows platforms
const cmdTimeout = 1500 * time.Millisecond
