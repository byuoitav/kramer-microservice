package videoswitcher

import (
	"errors"
	"fmt"
	"strings"

	"github.com/byuoitav/av-api/statusevaluators"
)

func GetCurrentInputByOutputPort(address, port string, readWelcome bool) (statusevaluators.Input, error) {
	command := fmt.Sprintf("#VID? %s", port)

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

		var i statusevaluators.Input
		i.Input = parts[0]
		return i, nil
	}

	logError(fmt.Sprintf("Incorrect response for command. (Response: %s)", resp))
	return statusevaluators.Input{}, errors.New(fmt.Sprintf("Incorrect response for command. (Response: %s)", resp))
}
