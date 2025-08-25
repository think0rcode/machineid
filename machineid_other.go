//go:build !linux && !darwin && !windows

package machineid

import (
	"context"
)

func getSMBIOSUUID(context.Context) (string, error) {
	return "", ErrUnsupported
}

func getInstallationID(context.Context) (string, error) {
	return "", ErrUnsupported
}
