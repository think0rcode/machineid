package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/think0rcode/machineid"
)

func main() {
	var logLevel = flag.String("log-level", "", "Set log level (error, info, debug). Can also use MACHINEID_LOG_LEVEL env var")
	var checkTools = flag.Bool("check-tools", true, "On Windows, check availability of PowerShell and WMIC")
	flag.Parse()

	fmt.Println("logLevel", *logLevel, "env", os.Getenv("MACHINEID_LOG_LEVEL"))

	// Configure logging based on flag or environment
	if *logLevel != "" {
		switch *logLevel {
		case "debug":
			machineid.SetLogLevel(machineid.LogLevelDebug)
		case "info":
			machineid.SetLogLevel(machineid.LogLevelInfo)
		case "error":
			machineid.SetLogLevel(machineid.LogLevelError)
		default:
			fmt.Fprintf(os.Stderr, "Invalid log level: %s. Valid values: error, info, debug\n", *logLevel)
			os.Exit(1)
		}
	} else {
		// Otherwise, use environment-based configuration (default)
		machineid.SetLogLevel(machineid.LogLevelDebug) // for debugging
	}

	slog.Info("starting machineid print tool")

	// Windows diagnostics (no-op on non-Windows builds)
	runWindowsDiagnostics(*checkTools)

	// Get raw components
	slog.Info("retrieving raw machine identifiers")
	bios, inst := machineid.RawID()
	fmt.Printf("bios=%s\ninst=%s\n", bios, inst)
	slog.Debug("raw identifiers retrieved", "bios", bios, "inst", inst)

	// Get hashed machine ID
	slog.Info("generating hashed machine ID")
	id, err := machineid.ID()
	if err != nil {
		slog.Error("failed to generate machine ID", "error", err)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("hashed=%s\n", id)
	slog.Info("machine ID tool completed successfully", "id", id)
}
