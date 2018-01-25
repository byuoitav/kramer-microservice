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
/*
func GetVolume(address string) bool {
	defer color.Unset()
	color.Set(color.FgYellow)
	ViaVolume := false

	log.Printf("Getting volume level of %s", address)

	var command ViaCommand
	resp, err := SendCommand(command, address)
	if err == nil && strings.Contai
	ns(resp, "Vol|Get|") {
		ViaVolume = true
	}

	return ViaVolume
}
*/
// GetVolume for a VIA device
func GetVolume(address string) (volcurrentlevel,error) {
	defer color.Unset()
	color.Set(color.FgYellow)

	var command ViaCommand
	command.Command = Vol
	command.Param1 = Get

	log.Printf("Sending command to get VIA Volume to %s", address)

	volcurrentlevel, err := SendCommand(command, address)
	if err != nil {
		return err
	}
  // Volume Get command in VIA API doesn't have any error handling so it only returns Vol|Get|XX or nothing
	if strings.Contains(resp, "Vol|GET|"){
		return volcurrentlevel
	}

	return errors.New(fmt.Sprintf("Incorrect response for command. (Response: %s)", resp))
}
