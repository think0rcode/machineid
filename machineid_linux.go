//go:build linux

package machineid

import (
	"context"
	"errors"
	"os"
	"strings"
)

func getSMBIOSUUID(_ context.Context) (string, error) {
	paths := []string{
		"/sys/class/dmi/id/product_uuid",
		"/sys/devices/virtual/dmi/id/product_uuid",
	}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		id := strings.TrimSpace(string(data))
		if id == "" || strings.EqualFold(id, "unknown") {
			continue
		}
		return strings.ToLower(id), nil
	}
	return "", errors.New("product uuid not available")
}

func getInstallationID(_ context.Context) (string, error) {
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
