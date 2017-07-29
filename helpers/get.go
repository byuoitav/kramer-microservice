package helpers

import (
	"errors"
	"fmt"
	"strings"
)

func GetCurrentInputByOutputPort(address string, port string) (string, error) {
	command := fmt.Sprintf("#VID? %s", port)
	resp, err := SendCommand(address, command)
	if err != nil {
		return "", err
	}

	if strings.Contains(resp, "VID") {
		parts := strings.Split(resp, "VID")
		resp = strings.TrimSpace(parts[1])

		parts = strings.Split(resp, ">")

		return parts[0], nil
	}
	return "", errors.New(fmt.Sprintf("Incorrect response for command. (Response: %s)", resp))
}
