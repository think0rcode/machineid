# machineid

Package `machineid` provides cross-platform generation of a stable machine
identifier for Go applications.

The identifier is built from three pieces of information:

- **SMBIOS hardware UUID** – a value provided by the system firmware. Hypervisors
  typically generate a new SMBIOS UUID for cloned virtual machines, ensuring that
  each clone receives a different identifier.
- **Operating system installation ID** – the OS's own identifier when available.
  Examples include `/etc/machine-id` on Linux and the `MachineGuid` registry value
  on Windows. This acts as an additional safeguard and fallback.
- **Operating system name** – included to distinguish IDs when a machine is
  re-imaged with a different OS.

The final ID is the SHA-256 hash of the string `os:<goos>|bios:<uuid>|inst:<id>`.
No MAC addresses or other volatile values are used, making the result stable even
on modern macOS versions that randomise network interfaces on each reboot.

## Usage

Install the module:

```bash
go get github.com/machineid/machineid
```

Then use it in your code:

```go
package main

import (
    "fmt"

    "github.com/machineid/machineid"
)

func main() {
    id, err := machineid.ID()
    if err != nil {
        panic(err)
    }
    fmt.Println(id)
}
```

For troubleshooting, `RawID` exposes the underlying components:

```go
bios, inst := machineid.RawID()
fmt.Println("bios:", bios, "install:", inst)
```

### CLI debugging tool

A small helper at `cmd/print` prints the raw components and the hashed ID. It also supports logging and, on Windows, optional tool checks:

```bash
# Print values with debug logs
go run ./cmd/print --log-level debug

# On Windows, also check PowerShell and WMIC availability
go run ./cmd/print --check-tools
```

## Supported platforms

- **Linux** – reads the SMBIOS UUID from `/sys/class/dmi/id/product_uuid`
  (falling back to `/sys/devices/virtual/dmi/id/product_uuid`) and the installation
  ID from `/etc/machine-id` or `/var/lib/dbus/machine-id`.
- **macOS** – obtains the hardware UUID from `ioreg` (falling back to
  `system_profiler`) and uses the system serial number as the installation ID,
  reusing the hardware UUID if unavailable.
- **Windows** – queries the SMBIOS UUID via PowerShell (prefers `Get-CimInstance`,
  falls back to `Get-WmiObject`, then `wmic`) and reads the installation ID
  from the `MachineGuid` registry value using the Go Windows registry API.

## Releasing

1. Commit any pending changes and update documentation.
2. Create a new semantic version tag:

   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```

   Go modules will pick up the latest tag automatically.

## License

This project is licensed under the terms of the Apache License 2.0. See [LICENSE](LICENSE).
