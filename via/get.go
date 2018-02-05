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
	ns(resp, "Vol|Get|") {volumefin
		ViaVolume = true
	}

	return ViaVolume
}
*/
// GetVolume for a VIA device
//func GetVolume(address string) (se.Volume, error) {
func GetVolume(address string) (int, error) {
	defer color.Unset()
	color.Set(color.FgYellow)

	var command ViaCommand
	command.Command = "Vol"
	command.Param1 = "Get"

	log.Printf("Sending command to get VIA Volume to %s", address)

	vollevel, err := SendCommand(command, address)

	if err != nil {
		return 0, err
	} else {
		//check the error first and then parse
		// Moved Parsing over to Common.go so anything can parse out the number
		/*
			re := regexp.MustCompile("[0-9]+")
			vol := re.FindString(volcurrentlevel)
			vfin, _ := strconv.Atoi(vol)
		*/

		// parse the returned string for the number that volume is set to
		volint, err := VolumeParse(vollevel)

		if err != nil {
			//return se.Volume{}, err
			return 0, err //passing 0 response along with error
		} else {
			// Volume Get command in VIA API doesn't have any error handling so it only returns Vol|Get|XX or nothing
			//if strings.Contains(volcurrentlevel, "Vol|GET|"){

			return volint, nil

			//}
		}
	}
}
