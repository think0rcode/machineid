//go:build linux

package machineid

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"strings"
)

func getSMBIOSUUID(_ context.Context) (string, error) {
	slog.Debug("getting SMBIOS UUID on Linux", "method", "sysfs")
	paths := []string{
		"/sys/class/dmi/id/product_uuid",
		"/sys/devices/virtual/dmi/id/product_uuid",
	}

	for _, p := range paths {
		slog.Debug("trying SMBIOS UUID path", "path", p)
		data, err := os.ReadFile(p)
		if err != nil {
			slog.Debug("failed to read SMBIOS UUID path", "path", p, "error", err)
			continue
		}
		id := strings.TrimSpace(string(data))
		if id == "" || strings.EqualFold(id, "unknown") {
			slog.Debug("invalid SMBIOS UUID found", "path", p, "value", id)
			continue
		}
		slog.Debug("SMBIOS UUID found", "path", p, "uuid", id)
		return strings.ToLower(id), nil
	}

	slog.Error("failed to find SMBIOS UUID on Linux", "paths_tried", paths)
	return "", errors.New("product uuid not available")
}

func getInstallationID(_ context.Context) (string, error) {
	slog.Debug("getting installation ID on Linux", "method", "machine-id")
	paths := []string{"/etc/machine-id", "/var/lib/dbus/machine-id"}

	for _, p := range paths {
		slog.Debug("trying installation ID path", "path", p)
		data, err := os.ReadFile(p)
		if err != nil {
			slog.Debug("failed to read installation ID path", "path", p, "error", err)
			continue
		}
		id := strings.TrimSpace(string(data))
		if id != "" {
			slog.Debug("installation ID found", "path", p, "id", id)
			return strings.ToLower(id), nil
		}
		slog.Debug("empty installation ID found", "path", p)
	}

	slog.Error("failed to find installation ID on Linux", "paths_tried", paths)
	return "", errors.New("machine-id not found")
}
