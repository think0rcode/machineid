# machineid

Package `machineid` provides cross-platform generation of a stable machine
identifier for Go applications.

The identifier is built from two components:

- **SMBIOS hardware UUID** – a value provided by the system firmware. Hypervisors
  typically generate a new SMBIOS UUID for cloned virtual machines, ensuring that
  each clone receives a different identifier.
- **Operating system installation ID** – the OS's own identifier when available.
  Examples include `/etc/machine-id` on Linux and the `MachineGuid` registry value
  on Windows. This acts as an additional safeguard and fallback.

The final ID is the SHA-256 hash of the two components concatenated together. No
MAC addresses or other volatile values are used, making the result stable even on
modern macOS versions that randomise network interfaces on each reboot.

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

## Supported platforms

- **Linux** – reads the SMBIOS UUID from `/sys/class/dmi/id/product_uuid` and the
  installation ID from `/etc/machine-id` (falling back to
  `/var/lib/dbus/machine-id`).
- **macOS** – obtains the SMBIOS UUID and serial number from `ioreg`.
- **Windows** – uses `wmic csproduct get UUID` for the SMBIOS UUID and the
  `MachineGuid` registry key for the installation ID.

## Releasing

1. Commit any pending changes and update documentation.
2. Create a new semantic version tag:

   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```

   Go modules will pick up the latest tag automatically.

## License

This project is licensed under the terms of the MIT license. See [LICENSE](LICENSE).
