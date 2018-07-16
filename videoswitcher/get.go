package videoswitcher

import (
	"errors"
	"fmt"
	"strings"

	"github.com/byuoitav/common/status"
)

func GetCurrentInputByOutputPort(address, port string, readWelcome bool) (status.Input, error) {
	command := fmt.Sprintf("#VID? %s", port)

	respChan := make(chan Response)

	c := CommandInfo{respChan, address, command, readWelcome}

	StartChannel <- c

	re := <-respChan

	resp := re.Response
	err := re.Err

	if err != nil {
		logError(err.Error())
		return status.Input{}, err
	}

	if strings.Contains(resp, "VID") {
		parts := strings.Split(resp, "VID")
		resp = strings.TrimSpace(parts[1])

		parts = strings.Split(resp, ">")

		var i status.Input
		i.Input = parts[0]
		return i, nil
	}

	logError(fmt.Sprintf("Incorrect response for command (%s). (Response: %s)", command, resp))
	return status.Input{}, errors.New(fmt.Sprintf("Incorrect response for command (%s). (Response: %s)", command, resp))
}
