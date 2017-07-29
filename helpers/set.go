package helpers

import (
	"errors"
	"fmt"
	"strings"
)

func SwitchInput(address string, input string, output string) error {
	command := fmt.Sprintf("#VID %s>%s", input, output)

	resp, err := SendCommand(address, command)
	if err != nil {
		return err
	}

	if strings.Contains(resp, "VID") {
		parts := strings.Split(resp, "VID")
		resp = strings.TrimSpace(parts[1])

		parts = strings.Split(resp, ">")

		if parts[0] == input && parts[1] == output {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Incorrect response for command. (Response: %s)", resp))
}
