//go:build windows

package machineid

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/sys/windows/registry"
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
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Cryptography`, registry.QUERY_VALUE|registry.WOW64_64KEY)
	if err != nil {
		return "", err
	}
	defer k.Close()
	guid, _, err := k.GetStringValue("MachineGuid")
	if err != nil {
		return "", errors.New("MachineGuid not found")
	}
	if guid == "" {
		return "", errors.New("machine guid empty")
	}
	return strings.ToLower(guid), nil
}
