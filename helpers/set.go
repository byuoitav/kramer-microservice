package helpers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func SwitchInput(address string, input int, output int) error {
	command := fmt.Sprintf("#VID %v>%v", input, output)

	resp, err := SendCommand(address, command)
	if err != nil {
		return err
	}

	if strings.Contains(resp, "VID") {
		parts := strings.Split(resp, "VID")
		resp = strings.TrimSpace(parts[1])

		parts = strings.Split(resp, ">")

		i, err := strconv.Atoi(parts[0])
		if err != nil {
			return err
		}
		o, err := strconv.Atoi(parts[1])
		if err != nil {
			return err
		}

		if i == input && o == output {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Incorrect response for command. (Response: %s)", resp))
}

func SetFrontLock(address string, state bool) error {
	var num int8
	if state {
		num = 1
	}

	command := fmt.Sprintf("#LOCK-FP %v", num)

	resp, err := SendCommand(address, command)
	if err != nil {
		return err
	}

	if strings.Contains(resp, "OK") {
		return nil
	}
	return errors.New(fmt.Sprintf("Incorrect response for command. (Response: %s)", resp))
}
