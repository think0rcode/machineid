package machineid

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
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

	bios, _ := getSMBIOSUUID(ctx)
	inst, _ := getInstallationID(ctx)

	bios = strings.ToLower(strings.TrimSpace(bios))
	inst = strings.ToLower(strings.TrimSpace(inst))

	if isZeroUUID(bios) {
		bios = ""
	}

	if bios == "" && inst == "" {
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
