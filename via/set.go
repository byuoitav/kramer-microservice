package via

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/fatih/color"
)

const REBOOT = "Reboot"
const RESET = "Reset"

func Reboot(address string) error {
	defer color.Unset()
	color.Set(color.FgYellow)

	var command ViaCommand
	command.Command = REBOOT

	log.Printf("Sending command %s to %s", REBOOT, address)

	_, err := SendCommand(command, address)
	if err != nil {
		return err
	}

	return nil
}

func Reset(address string) error {
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
