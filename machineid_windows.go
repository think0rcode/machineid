//go:build windows

package machineid

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func getSMBIOSUUID(ctx context.Context) (string, error) {
	// Try modern PowerShell CIM first (PS 3.0+; Server 2012+)
	slog.Debug("getting SMBIOS UUID on Windows", "method", "powershell_cim")
	out := run(ctx, "powershell", "-NoProfile", "-Command", "(Get-CimInstance Win32_ComputerSystemProduct).UUID")

	slog.Debug("powershell CIM output", "output", out)

	if u := firstUUID(out); u != "" {
		slog.Debug("SMBIOS UUID found via PowerShell CIM", "uuid", u)
		return u, nil
	}

	// Fall back to older WMI cmdlet (available in PS 2.0+)
	slog.Debug("CIM failed, trying PowerShell WMI", "method", "powershell_wmi")
	out = run(ctx, "powershell", "-NoProfile", "-Command", "(Get-WmiObject Win32_ComputerSystemProduct).UUID")
	slog.Debug("powershell WMI output", "output", out)
	if u := firstUUID(out); u != "" {
		slog.Debug("SMBIOS UUID found via PowerShell WMI", "uuid", u)
		return u, nil
	}

	// Final fallback to legacy WMIC (deprecated/removed on newer Windows)
	slog.Debug("PowerShell failed, trying WMIC", "method", "wmic")
	out = run(ctx, "wmic", "csproduct", "get", "uuid")
	slog.Debug("wmic output", "output", out)
	if u := firstUUID(out); u != "" {
		slog.Debug("SMBIOS UUID found via WMIC", "uuid", u)
		return u, nil
	}

	slog.Error("failed to find SMBIOS UUID on Windows", "methods_tried", []string{"powershell_cim", "powershell_wmi", "wmic"})
	return "", errors.New("no uuid found")
}

func getInstallationID(ctx context.Context) (string, error) {
	slog.Debug("getting installation ID on Windows", "method", "registry_machine_guid")
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Cryptography`, registry.QUERY_VALUE|registry.WOW64_64KEY)
	if err != nil {
		slog.Error("failed to open registry key for installation ID", "key", `SOFTWARE\Microsoft\Cryptography`, "error", err)
		return "", err
	}
	defer k.Close()

	// dump all key and values in the key
	values, err := k.ReadValueNames(0)
	if err != nil {
		slog.Error("failed to read values from registry", "error", err)
		return "", err
	}
	slog.Debug("registry key values", "values", values)

	guid, _, err := k.GetStringValue("MachineGuid")
	if err != nil {
		slog.Error("failed to read MachineGuid from registry", "error", err)
		return "", errors.New("MachineGuid not found")
	}
	if guid == "" {
		slog.Error("MachineGuid is empty", "value", guid)
		return "", errors.New("machine guid empty")
	}

	slog.Debug("installation ID found via registry", "guid", guid)
	return strings.ToLower(guid), nil
}
