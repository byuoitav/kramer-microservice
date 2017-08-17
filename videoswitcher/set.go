package videoswitcher

import (
	"errors"
	"fmt"
	"strings"
)

func SwitchInput(address, input, output string, readWelcome bool) error {
	command := fmt.Sprintf("#VID %s>%s", input, output)

	resp, err := SendCommand(address, command, readWelcome)
	if err != nil {
		logError(err.Error())
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

	logError(fmt.Sprintf("Incorrect response for command. (Response: %s)", resp))
	return errors.New(fmt.Sprintf("Incorrect response for command. (Response: %s)", resp))
}

func SetFrontLock(address string, state, readWelcome bool) error {
	var num int8
	if state {
		num = 1
	}

	command := fmt.Sprintf("#LOCK-FP %v", num)

	resp, err := SendCommand(address, command, readWelcome)
	if err != nil {
		logError(err.Error())
		return err
	}

	if strings.Contains(resp, "OK") {
		return nil
	}

	logError(fmt.Sprintf("Incorrect response for command. (Response: %s)", resp))
	return errors.New(fmt.Sprintf("Incorrect response for command. (Response: %s)", resp))
}
