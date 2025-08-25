package machineid

import (
	"bytes"
	"context"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const cmdTimeout = 1500 * time.Millisecond

func run(ctx context.Context, name string, args ...string) string {
	cmd := exec.CommandContext(ctx, name, args...)
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	_ = cmd.Run()
	return b.String()
}

var uuidRE = regexp.MustCompile(`(?i)[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)

func firstUUID(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	m := uuidRE.FindString(s)
	return strings.ToLower(m)
}

func isZeroUUID(u string) bool {
	u = strings.ToLower(strings.TrimSpace(u))
	zero := "00000000-0000-0000-0000-000000000000"
	ffff := "ffffffff-ffff-ffff-ffff-ffffffffffff"
	return u == zero || u == ffff
}
