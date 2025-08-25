// errors.go
package machineid

import "errors"

// ErrUnsupported indicates the current platform is not supported
// by the SMBIOS/installation ID backends.
var ErrUnsupported = errors.New("machineid: unsupported platform")
