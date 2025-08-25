# Logging Documentation

This document describes the logging functionality added to the machineid library.

## Overview

The machineid library now includes comprehensive structured logging using Go's standard `log/slog` package. The logging provides detailed insights into the machine ID generation process across all supported platforms.

## Log Levels

The library supports three log levels:

- **Error**: Shows only error messages (default for production use)
- **Info**: Shows informational messages about the main operations
- **Debug**: Shows detailed debug information including command execution and intermediate values

## Configuration

### Environment Variable

Set the `MACHINEID_LOG_LEVEL` environment variable:

```bash
export MACHINEID_LOG_LEVEL=debug    # Show all logs
export MACHINEID_LOG_LEVEL=info     # Show info and error logs
export MACHINEID_LOG_LEVEL=error    # Show only error logs (default)
```

### Programmatic Configuration

```go
import "github.com/think0rcode/machineid"

// Set log level programmatically
machineid.SetLogLevel(machineid.LogLevelDebug)
machineid.SetLogLevel(machineid.LogLevelInfo)
machineid.SetLogLevel(machineid.LogLevelError)

// Configure from environment variable
machineid.SetLogLevelFromEnv()
```

### Command Line Tool

The `cmd/print` tool supports a `--log-level` flag:

```bash
go run ./cmd/print --log-level debug
go run ./cmd/print --log-level info
go run ./cmd/print --log-level error
```

## What Gets Logged

### Error Level
- Platform not supported errors
- Failed to read identifiers
- Registry access failures (Windows)
- File system access failures (Linux)

### Info Level
- Machine ID generation start/completion
- Successful identifier retrieval
- Tool execution start/completion

### Debug Level
- Individual command executions with arguments and output
- Platform-specific method attempts
- UUID extraction and validation
- Hash generation process
- Raw identifier values (useful for diagnostics)

## Platform-Specific Logging

### Darwin (macOS)
- Logs `ioreg` and `system_profiler` command executions
- Shows serial number and hardware UUID retrieval
- Fallback mechanism logging

### Linux
- Logs file system reads from `/sys/class/dmi/id/` and `/etc/machine-id`
- Shows which paths were tried and which succeeded
- Invalid value detection (empty or "unknown")

### Windows
- Logs PowerShell and WMI command executions
- Registry key access for MachineGuid
- Shows method fallback from PowerShell to WMI

### Other Platforms
- Logs unsupported platform warnings

## Example Debug Output

```
time=2025-08-25T13:41:02.753+08:00 level=DEBUG msg="generating machine ID" os=darwin
time=2025-08-25T13:41:02.753+08:00 level=DEBUG msg="getting SMBIOS UUID on Darwin" method=ioreg
time=2025-08-25T13:41:02.753+08:00 level=DEBUG msg="running command" command=ioreg args="[-rd1 -c IOPlatformExpertDevice]"
time=2025-08-25T13:41:02.768+08:00 level=DEBUG msg="command completed successfully" command=ioreg args="[-rd1 -c IOPlatformExpertDevice]" output_length=1878
time=2025-08-25T13:41:02.769+08:00 level=DEBUG msg="extracted UUID from string" input_length=1871 uuid_found=true uuid=9129daec-8350-51f5-9247-f1db82392c45
time=2025-08-25T13:41:02.769+08:00 level=DEBUG msg="SMBIOS UUID found via ioreg" uuid=9129daec-8350-51f5-9247-f1db82392c45
```

## Performance Considerations

- Error level logging has minimal performance impact
- Debug level logging includes command output and may impact performance in high-frequency scenarios
- All logging uses structured logging with key-value pairs for easy parsing
- Logging is disabled by default (error level) to avoid noise in production environments

## Integration Notes

- The logging system is initialized automatically when the package is imported
- Default configuration reads from the `MACHINEID_LOG_LEVEL` environment variable
- All log output goes to `stderr` to avoid interfering with program output
- Logging configuration is global and affects all uses of the library within a process
