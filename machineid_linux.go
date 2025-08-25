//go:build linux

package machineid

import (
	"errors"
	"os"
	"strings"
)

func getSMBIOSUUID() (string, error) {
	data, err := os.ReadFile("/sys/class/dmi/id/product_uuid")
	if err != nil {
		return "", err
	}
	id := strings.TrimSpace(string(data))
	if id == "" || strings.EqualFold(id, "unknown") {
		return "", errors.New("product uuid not available")
	}
	return strings.ToLower(id), nil
}

func getInstallationID() (string, error) {
	paths := []string{"/etc/machine-id", "/var/lib/dbus/machine-id"}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		id := strings.TrimSpace(string(data))
		if id != "" {
			return strings.ToLower(id), nil
		}
	}
	return "", errors.New("machine-id not found")
}
