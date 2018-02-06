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

func SetVolume(address string, volumec string) (string, error) {
	defer color.Unset()
	color.Set(color.FgYellow)

	var command ViaCommand
	command.Command = "Vol"
	command.Param1 = "Set"
	command.Param2 = volumec

	log.Printf("Sending volume set command to %s", address)

	resp, err := SendCommand(command, address)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error in setting volume on %s", address))
	}

	// Error handling - Handle with care
	// if command returns vol|set|Error1 - value is outside the bounds of 0-100
	// if command returns vol|set|Error2 - no value set for volume
	if strings.Contains(resp, "Error1") {
		log.Printf("Volume command error - volume value %s is outside the bounds of 0-100", volumec)
		return "", errors.New(fmt.Sprintf("Volume value is outside the bounds of 0-100"))
	} else if strings.Contains(resp, "Error2") {
		log.Printf("Volume command error - volume value was not in the command passed to %s", address)
		return "", errors.New(fmt.Sprintf("Volume value was not in the command passed"))
	}

	return resp, nil

}
