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

	var command Command
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

	var command Command
	command.Command = "Vol"
	command.Param1 = "Get"

	log.Printf("Sending command to get VIA Volume to %s", address)
	// Note: Volume Get command in VIA API doesn't have any error handling so it only returns Vol|Get|XX or nothing
	// I am still checking for errors just in case something else fails during execution
	vollevel, _ := SendCommand(command, address)

	return VolumeParse(vollevel)
}
