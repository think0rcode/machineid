//go:build darwin

package machineid

import (
	"fmt"
	"os/exec"
	"strings"
)

func getIOPlatformProperty(prop string) (string, error) {
	out, err := exec.Command("ioreg", "-rd1", "-c", "IOPlatformExpertDevice").Output()
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "\""+prop+"\"") {
			parts := strings.Split(line, "=")
			if len(parts) != 2 {
				continue
			}
			v := strings.Trim(parts[1], " \"")
			if v != "" {
				return strings.ToLower(v), nil
			}
		}
	}
	return "", fmt.Errorf("%s not found", prop)
}

func getSMBIOSUUID() (string, error) {
	return getIOPlatformProperty("IOPlatformUUID")
}

func getInstallationID() (string, error) {
	// macOS does not expose an explicit installation ID. We approximate this
	// using the serial number from the IORegistry, which changes on hardware
	// replacement and differs across VM clones.
	return getIOPlatformProperty("IOPlatformSerialNumber")
}
