package videoswitcher

import (
	"errors"
	"fmt"
	"strings"

	"github.com/byuoitav/av-api/statusevaluators"
)

func SwitchInput(address, input, output string, readWelcome bool) (statusevaluators.Input, error) {
	command := fmt.Sprintf("#VID %s>%s", input, output)

	respChan := make(chan Response)
	c := CommandInfo{respChan, address, command, readWelcome}

	StartChannel <- c

	re := <-respChan

	resp := re.Response
	err := re.Err

	if err != nil {
		logError(err.Error())
		return statusevaluators.Input{}, err
	}

	if strings.Contains(resp, "VID") {
		parts := strings.Split(resp, "VID")
		resp = strings.TrimSpace(parts[1])

		parts = strings.Split(resp, ">")

		if parts[0] == input && parts[1] == output {
			var i statusevaluators.Input
			i.Input = input
			return i, err
		}
	}

	logError(fmt.Sprintf("Incorrect response for command. (Response: %s)", resp))
	return statusevaluators.Input{}, errors.New(fmt.Sprintf("Incorrect response for command. (Response: %s)", resp))
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

	logError(fmt.Sprintf("Incorrect response for command. (Response: %s)", resp))
	return errors.New(fmt.Sprintf("Incorrect response for command. (Response: %s)", resp))
}
