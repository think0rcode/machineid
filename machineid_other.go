//go:build !linux && !darwin && !windows

package machineid

import (
	"context"
	"errors"
)

func getSMBIOSUUID(context.Context) (string, error) {
	return "", errors.New("unsupported platform")
}

func getInstallationID(context.Context) (string, error) {
	return "", errors.New("unsupported platform")
}
