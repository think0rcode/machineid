//go:build !linux && !darwin && !windows

package machineid

import (
	"context"

	"log/slog"
)

func getSMBIOSUUID(context.Context) (string, error) {
	slog.Error("SMBIOS UUID not supported on this platform")
	return "", ErrUnsupported
}

func getInstallationID(context.Context) (string, error) {
	slog.Error("installation ID not supported on this platform")
	return "", ErrUnsupported
}
