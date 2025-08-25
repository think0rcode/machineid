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
	var checkTools = flag.Bool("check-tools", false, "On Windows, check availability of PowerShell and WMIC")
	flag.Parse()

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
	}
	// Otherwise, use environment-based configuration (default)

	slog.Info("starting machineid print tool")

	// Windows diagnostics (no-op on non-Windows builds)
	runWindowsDiagnostics(*checkTools)

	// Get raw components
	slog.Info("retrieving raw machine identifiers")
	bios, inst, biosErr, instErr := machineid.RawID()
	if biosErr != nil {
		slog.Error("failed to retrieve BIOS UUID", "error", biosErr)
	}
	if instErr != nil {
		slog.Error("failed to retrieve installation ID", "error", instErr)
	}
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

	// want the process did not quit, until i press enter
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}
