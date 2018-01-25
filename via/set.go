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

func SetVolume(address string,ViaVolumeLevel int) error {
	defer color.Unset()
	color.Set(color.FgYellow)

	var command ViaCommand
	command.Command = Vol
	command.Param1 = Set
	command.Param2 = ViaVolumeLevel

	log.Printf("Sending volume set command to %s", address)

	resp, err := SendCommand(command, address)
	if err != nil {
		return err
	}
  // Error handling - Handle with care
	// Error1 - value is outside the bounds of 0-100
	// Error2 - no value set for volume
	if strings.Contains(resp, "Error1") {
		return resp
		log.Printf("Volume command error - volume value %s is outside the bounds of 0-100", ViaVolumeLevel)
	} else if strings.Contains(resp, "Error2"){
		return resp
		log.Printf("Volume command error - volume value was not in the command passed to %s", address)
	} else {
		//r, _ := regexp.Compile("/\|\d/g")
    r, _ := regexp.Compile("/\d/g")
		VolumeSetFin := r.MatchString(resp) //matching strings or should we convert to integer
	}
  if VolumeSetFin != ViaVolumeLevel{
		log.Printf("Volume command error - volume did not change to %s as requested", ViaVolumeLevel)
	}else{
		log.Printf("Volume change successful")
	}
	return errors.New(fmt.Sprintf("Incorrect response for command. (Response: %s)", resp))
}
