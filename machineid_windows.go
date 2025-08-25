//go:build windows

package machineid

import (
	"errors"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func getSMBIOSUUID() (string, error) {
	out, err := exec.Command("wmic", "csproduct", "get", "UUID").CombinedOutput()
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.EqualFold(line, "uuid") {
			continue
		}
		return strings.ToLower(line), nil
	}
	return "", errors.New("wmic returned no uuid")
}

func getInstallationID() (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Cryptography`, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()
	s, _, err := k.GetStringValue("MachineGuid")
	if err != nil {
		return "", err
	}
	if s == "" {
		return "", errors.New("machine guid empty")
	}
	return strings.ToLower(s), nil
}
