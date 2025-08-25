//go:build windows

package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func runSilent(timeout time.Duration, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out

	done := make(chan error, 1)
	go func() { done <- cmd.Run() }()

	select {
	case err := <-done:
		return out.String(), err
	case <-time.After(timeout):
		_ = cmd.Process.Kill()
		return "", fmt.Errorf("timeout running %s", name)
	}
}

// runWindowsDiagnostics prints basic environment diagnostics for Windows.
// It mirrors the checks from cmd/wincheck, and is a no-op on other OSes.
func runWindowsDiagnostics(checkTools bool) {
	if !checkTools {
		return
	}

	fmt.Printf("wincheck: GOOS=%s GOARCH=%s\n", runtime.GOOS, runtime.GOARCH)

	// PowerShell version
	psOut, psErr := runSilent(2*time.Second, "powershell", "-NoProfile", "-Command", "[string]$PSVersionTable.PSVersion")
	psOut = strings.TrimSpace(psOut)
	if psErr != nil || psOut == "" {
		fmt.Printf("powershell: unavailable (%v)\n", psErr)
	} else {
		fmt.Printf("powershell: %s\n", strings.ReplaceAll(psOut, "\n", " "))
	}

	// WMIC availability (deprecated on newer Windows)
	if _, err := runSilent(1*time.Second, "wmic", "/?"); err != nil {
		fmt.Printf("wmic: unavailable (%v)\n", err)
	} else {
		fmt.Println("wmic: available")
	}
}
