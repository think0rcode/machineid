package machineid

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
)

// ID returns a hashed identifier for the current machine. It combines the
// operating system name, SMBIOS hardware UUID, and an operating system
// installation identifier when available. The identifier is stable for a
// given machine yet changes when the VM or host is cloned.
func ID() (string, error) {
	slog.Debug("generating machine ID", "os", runtime.GOOS)

	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()

	bios, biosErr := getSMBIOSUUID(ctx)
	inst, instErr := getInstallationID(ctx)

	slog.Debug("retrieved raw identifiers",
		"bios_uuid", bios,
		"bios_error", biosErr,
		"installation_id", inst,
		"installation_error", instErr)

	bios = strings.ToLower(strings.TrimSpace(bios))
	inst = strings.ToLower(strings.TrimSpace(inst))

	if isZeroUUID(bios) {
		bios = ""
	}

	if bios == "" && inst == "" {
		slog.Error("failed to retrieve any machine identifiers", "bios_empty", true, "installation_empty", true)
		return "", errors.New("machineid: unable to read BIOS UUID or installation ID")
	}

	if bios == "" && inst == "" {
		// If both backends are unsupported, surface a more specific error
		if biosErr != nil && instErr != nil &&
			strings.Contains(biosErr.Error(), "unsupported platform") &&
			strings.Contains(instErr.Error(), "unsupported platform") {

			slog.Error("platform not supported", "bios_error", biosErr, "installation_error", instErr)
			return "", errors.New("machineid: unsupported platform")
		}

		// Otherwise, return a combined error to aid debugging
		if biosErr != nil && instErr != nil {
			slog.Error("failed to read both identifiers", "bios_error", biosErr, "installation_error", instErr)
			return "", fmt.Errorf("machineid: unable to read BIOS UUID (%w) or installation ID (%w)", biosErr, instErr)
		} else if biosErr != nil {
			slog.Error("failed to read BIOS UUID", "bios_error", biosErr, "installation_id_empty", true)
			return "", fmt.Errorf("machineid: unable to read BIOS UUID (%w), installation ID empty", biosErr)
		} else if instErr != nil {
			slog.Error("failed to read installation ID", "installation_error", instErr, "bios_uuid_empty", true)
			return "", fmt.Errorf("machineid: unable to read installation ID (%w), BIOS UUID empty", instErr)
		}

		return "", errors.New("machineid: unable to read BIOS UUID or installation ID")
	}

	slog.Debug("generating hash from identifiers",
		"os", runtime.GOOS,
		"bios_uuid", bios,
		"installation_id", inst)

	h := sha256.New()
	hashInput := "os:" + runtime.GOOS + "|bios:" + bios + "|inst:" + inst
	h.Write([]byte(hashInput))
	result := hex.EncodeToString(h.Sum(nil))

	slog.Info("machine ID generated successfully", "id", result)
	return result, nil
}

// RawID returns the unprocessed SMBIOS UUID and installation ID used to
// build the hashed identifier. It is primarily useful for diagnostics.
func RawID() (biosUUID, installID string) {
	slog.Debug("retrieving raw machine identifiers for diagnostics")

	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()
	bios, biosErr := getSMBIOSUUID(ctx)
	inst, instErr := getInstallationID(ctx)

	biosResult := strings.TrimSpace(bios)
	instResult := strings.TrimSpace(inst)

	slog.Debug("raw identifiers retrieved",
		"bios_uuid", biosResult,
		"bios_error", biosErr,
		"installation_id", instResult,
		"installation_error", instErr)

	return biosResult, instResult
}
