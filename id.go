package machineid

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// ID returns a hashed identifier for the current machine. It combines the
// operating system name, SMBIOS hardware UUID, and an operating system
// installation identifier when available. The identifier is stable for a
// given machine yet changes when the VM or host is cloned.
func ID() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()

	bios, biosErr := getSMBIOSUUID(ctx)
	inst, instErr := getInstallationID(ctx)

	bios = strings.ToLower(strings.TrimSpace(bios))
	inst = strings.ToLower(strings.TrimSpace(inst))

	if isZeroUUID(bios) {
		bios = ""
	}

	if bios == "" && inst == "" {
		return "", errors.New("machineid: unable to read BIOS UUID or installation ID")
	}

	if bios == "" && inst == "" {
		// If both backends are unsupported, surface a more specific error
		if biosErr != nil && instErr != nil &&
			strings.Contains(biosErr.Error(), "unsupported platform") &&
			strings.Contains(instErr.Error(), "unsupported platform") {
			return "", errors.New("machineid: unsupported platform")
		}

		// Otherwise, return a combined error to aid debugging
		if biosErr != nil && instErr != nil {
			return "", fmt.Errorf("machineid: unable to read BIOS UUID (%w) or installation ID (%w)", biosErr, instErr)
		} else if biosErr != nil {
			return "", fmt.Errorf("machineid: unable to read BIOS UUID (%w), installation ID empty", biosErr)
		} else if instErr != nil {
			return "", fmt.Errorf("machineid: unable to read installation ID (%w), BIOS UUID empty", instErr)
		}

		return "", errors.New("machineid: unable to read BIOS UUID or installation ID")
	}

	h := sha256.New()
	h.Write([]byte("os:" + runtime.GOOS + "|bios:" + bios + "|inst:" + inst))
	return hex.EncodeToString(h.Sum(nil)), nil
}

// RawID returns the unprocessed SMBIOS UUID and installation ID used to
// build the hashed identifier. It is primarily useful for diagnostics.
func RawID() (biosUUID, installID string) {
	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()
	bios, _ := getSMBIOSUUID(ctx)
	inst, _ := getInstallationID(ctx)
	return strings.TrimSpace(bios), strings.TrimSpace(inst)
}
