package machineid

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// ID returns a hashed identifier for the current machine. It combines the
// SMBIOS hardware UUID with an operating system installation identifier when
// available and hashes the result using SHA-256. The identifier is stable for a
// given machine yet changes when the VM or host is cloned by virtue of the
// SMBIOS UUID changing.
func ID() (string, error) {
	hw, errHW := getSMBIOSUUID()
	osID, errOS := getInstallationID()
	if errHW != nil && errOS != nil {
		return "", fmt.Errorf("cannot obtain identifiers: %v; %v", errHW, errOS)
	}
	sum := sha256.Sum256([]byte(hw + "|" + osID))
	return hex.EncodeToString(sum[:]), nil
}
