//go:build windows

package machineid

import "time"

// Windows can be slower to spin up PowerShell and CIM providers,
// so use a slightly longer timeout here.
const cmdTimeout = 10 * time.Second
