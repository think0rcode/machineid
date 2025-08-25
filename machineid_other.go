//go:build !linux && !darwin && !windows

package machineid

import "errors"

func getSMBIOSUUID() (string, error) {
	return "", errors.New("unsupported platform")
}

func getInstallationID() (string, error) {
	return "", errors.New("unsupported platform")
}
