package via

import (
	"log"
	"strings"

	"github.com/fatih/color"
)

func IsConnected(address string) bool {
	defer color.Unset()
	color.Set(color.FgYellow)
	connected := false

	log.Printf("Getting connected status of %s", address)

	var command ViaCommand
	resp, err := SendCommand(command, address)
	if err == nil && strings.Contains(resp, "Successful") {
		connected = true
	}

	return connected
}

// GetVolume for a VIA device
func GetVolume(address string) (int, error) {

	defer color.Unset()
	color.Set(color.FgYellow)

	var command ViaCommand
	command.Command = "Vol"
	command.Param1 = "Get"

	log.Printf("Sending command to get VIA Volume to %s", address)
	// Note: Volume Get command in VIA API doesn't have any error handling so it only returns Vol|Get|XX or nothing
	// I am still checking for errors just in case something else fails during execution
	vollevel, err := SendCommand(command, address)
	//check the error first and then parse
	if err != nil {
		return 0, err
	} else {
		// parse the returned string for the number that volume is set to
		// code in common.go so other parts can use the same parser function
		volint, err := VolumeParse(vollevel)

		if err != nil {
			return 0, err //passing 0 response along with error
		} else {
			return volint, nil
		}
	}
}
