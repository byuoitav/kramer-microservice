package helpers

import (
	"fmt"
)

func SwitchInput(address string, input string, output string) error {
	command := fmt.Sprintf("#VID %s>%s", input, output)

	_, err := SendCommand(address, command)
	if err != nil {
		return err
	}

	return nil
}
