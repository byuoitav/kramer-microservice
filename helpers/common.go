package helpers

import (
	"log"

	"github.com/fatih/color"
)

// takes a command and sends it to the device
// should it also return what the response of the query was?
func SendCommand(address string, command string) (string, error) {
	defer color.Unset()
	color.Set(color.FgCyan)

	// open telnet connection with address
	log.Printf("Opening telnet connection with %s", address)

	// write command
	log.Printf("Sending command '%s'", command)

	// get response
	var response string
	log.Printf("Response from device: %s", response)

	return response, nil
}
