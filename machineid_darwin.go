//go:build darwin

package machineid

import (
	"context"
	"errors"
	"strings"
)

func getSMBIOSUUID(ctx context.Context) (string, error) {
	out := run(ctx, "ioreg", "-rd1", "-c", "IOPlatformExpertDevice")
	if u := firstUUID(out); u != "" {
		return u, nil
	}
	out = run(ctx, "system_profiler", "SPHardwareDataType")
	if u := firstUUID(out); u != "" {
		return u, nil
	}
	return "", errors.New("IOPlatformUUID not found")
}

func getInstallationID(ctx context.Context) (string, error) {
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
				return strings.ToLower(v), nil
			}
		}
	}
	// Fallback: reuse hardware UUID if serial number unavailable
	return getSMBIOSUUID(ctx)
}
