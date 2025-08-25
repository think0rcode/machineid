package machineid

import (
	"bytes"
	"context"
	"log/slog"
	"os/exec"
	"regexp"
	"strings"
)

func run(ctx context.Context, name string, args ...string) string {
	slog.Debug("running command", "command", name, "args", args)
	cmd := exec.CommandContext(ctx, name, args...)
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b

	if err := cmd.Run(); err != nil {
		slog.Debug("command failed", "command", name, "args", args, "error", err, "output", b.String())
		return ""
	}

	output := b.String()
	slog.Debug("command completed successfully", "command", name, "args", args, "output_length", len(output))
	return output
}

var uuidRE = regexp.MustCompile(`(?i)[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)

func firstUUID(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		slog.Debug("firstUUID called with empty string")
		return ""
	}
	m := uuidRE.FindString(s)
	slog.Debug("extracted UUID from string", "input_length", len(s), "uuid_found", m != "", "uuid_redacted", redactID(m))
	return strings.ToLower(m)
}

func isZeroUUID(u string) bool {
	u = strings.ToLower(strings.TrimSpace(u))
	zero := "00000000-0000-0000-0000-000000000000"
	ffff := "ffffffff-ffff-ffff-ffff-ffffffffffff"
	isZero := u == zero || u == ffff
	slog.Debug("checking if UUID is zero", "uuid_redacted", redactID(u), "is_zero", isZero)
	return isZero
}

// redactID returns a privacy-preserving representation like "***abcd" (last 4).
func redactID(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if len(s) <= 4 {
		return "***"
	}
	return "***" + s[len(s)-4:]
}
