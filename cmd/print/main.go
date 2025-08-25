package main

import (
	"fmt"
	"os"

	"github.com/think0rcode/machineid"
)

func main() {
	// Get raw components
	bios, inst := machineid.RawID()
	fmt.Printf("bios=%s\ninst=%s\n", bios, inst)

	// Get hashed machine ID
	id, err := machineid.ID()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("hashed=%s\n", id)
}
