//go:build !windows

package machineid

import "time"

// Default timeout for non-Windows platforms
const cmdTimeout = 5 * time.Second
