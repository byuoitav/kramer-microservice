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
func GetVolume(address string) error {
	defer color.Unset()
	color.Set(color.FgYellow)

	var command ViaCommand
	command.Command = RESET

	log.Printf("Sending command %s to %s", RESET, address)

	resp, err := SendCommand(command, address)
	if err != nil {
		return err
	}

	if strings.Contains(resp, RESET) && strings.Contains(resp, "1") {
		return nil
	}

	return errors.New(fmt.Sprintf("Incorrect response for command. (Response: %s)", resp))
}
