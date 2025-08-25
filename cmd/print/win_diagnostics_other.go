//go:build !windows

package main

// runWindowsDiagnostics is a no-op on non-Windows platforms.
func runWindowsDiagnostics(_ bool) {}
