package helpers

import (
	"fmt"
	"log"
)

func SwitchInput(address string, input string, output string) error {
	command := fmt.Sprintf("#VID %s>%s", input, output)

	err := SendCommand(address, command)
	if err != nil {
		return err
	}

	return nil
}

// takes a command and sends it to the device
// should it also return what the response of the query was?
func SendCommand(address string, command string) error {
	// open telnet connection with address

	// append # to command
	log.Printf("Sending command '%s' to %s", command, address)

	// write command

	return nil
}
