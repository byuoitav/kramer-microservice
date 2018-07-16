package videoswitcher

import (
	"errors"
	"fmt"
	"strings"

	"github.com/byuoitav/common/status"
)

func SwitchInput(address, input, output string, readWelcome bool) (status.Input, error) {
	command := fmt.Sprintf("#VID %s>%s", input, output)

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

		if parts[0] == input && parts[1] == output {
			var i status.Input
			i.Input = input
			return i, err
		}
	}

	logError(fmt.Sprintf("Incorrect response for command (%s). (Response: %s)", command, resp))
	return status.Input{}, errors.New(fmt.Sprintf("Incorrect response for command (%s). (Response: %s)", command, resp))
}

func SetFrontLock(address string, state, readWelcome bool) error {
	var num int8
	if state {
		num = 1
	}
	command := fmt.Sprintf("#LOCK-FP %v", num)

	respChan := make(chan Response)
	c := CommandInfo{respChan, address, command, readWelcome}

	StartChannel <- c

	re := <-respChan

	resp := re.Response
	err := re.Err

	if err != nil {
		logError(err.Error())
		return err
	}

	if strings.Contains(resp, "OK") {
		return nil
	}

	logError(fmt.Sprintf("Incorrect response for command (%s). (Response: %s)", command, resp))
	return errors.New(fmt.Sprintf("Incorrect response for command (%s). (Response: %s)", command, resp))
}
