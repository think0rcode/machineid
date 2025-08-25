//go:build windows

package machineid

import (
	"context"
	"errors"
	"strings"
)

func getSMBIOSUUID(ctx context.Context) (string, error) {
	out := run(ctx, "powershell", "-NoProfile", "-Command", "(Get-CimInstance Win32_ComputerSystemProduct).UUID")
	if u := firstUUID(out); u != "" {
		return u, nil
	}
	out = run(ctx, "wmic", "csproduct", "get", "uuid")
	if u := firstUUID(out); u != "" {
		return u, nil
	}
	return "", errors.New("no uuid found")
}

func getInstallationID(ctx context.Context) (string, error) {
	out := run(ctx, "cmd", "/C", `reg query HKLM\SOFTWARE\Microsoft\Cryptography /v MachineGuid`)
	lines := strings.Split(out, "\n")
	for _, ln := range lines {
		if strings.Contains(ln, "MachineGuid") {
			fields := strings.Fields(ln)
			if len(fields) > 0 {
				return strings.ToLower(fields[len(fields)-1]), nil
			}
		}
	}
	return "", errors.New("MachineGuid not found")
}
