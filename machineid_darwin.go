//go:build darwin

package machineid

import (
	"context"
	"errors"
	"log/slog"
	"strings"
)

func getSMBIOSUUID(ctx context.Context) (string, error) {
	slog.Debug("getting SMBIOS UUID on Darwin", "method", "ioreg")
	out := run(ctx, "ioreg", "-rd1", "-c", "IOPlatformExpertDevice")
	if u := firstUUID(out); u != "" {
		slog.Debug("SMBIOS UUID found via ioreg", "uuid", u)
		return u, nil
	}

	slog.Debug("ioreg failed, trying system_profiler", "method", "system_profiler")
	out = run(ctx, "system_profiler", "SPHardwareDataType")
	if u := firstUUID(out); u != "" {
		slog.Debug("SMBIOS UUID found via system_profiler", "uuid", u)
		return u, nil
	}

	slog.Error("failed to find SMBIOS UUID on Darwin", "methods_tried", []string{"ioreg", "system_profiler"})
	return "", errors.New("IOPlatformUUID not found")
}

func getInstallationID(ctx context.Context) (string, error) {
	slog.Debug("getting installation ID on Darwin", "method", "serial_number")
	out := run(ctx, "ioreg", "-rd1", "-c", "IOPlatformExpertDevice")
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "\"IOPlatformSerialNumber\"") {
			parts := strings.Split(line, "=")
			if len(parts) != 2 {
				break
			}
			v := strings.Trim(parts[1], " \"")
			if v != "" {
				slog.Debug("installation ID found via serial number", "serial", v)
				return strings.ToLower(v), nil
			}
		}
	}

	slog.Debug("serial number not found, falling back to hardware UUID")
	// Fallback: reuse hardware UUID if serial number unavailable

	uuid, err := getSMBIOSUUID(ctx)
	if err != nil {
		slog.Error("failed to get installation ID on Darwin", "serial_found", false, "uuid_error", err)
		return "", errors.New("neither serial number nor hardware UUID found")
	}

	slog.Debug("installation ID using hardware UUID fallback", "uuid", uuid)
	return uuid, nil
}
