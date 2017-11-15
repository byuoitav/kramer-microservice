package p2000

import "log"

func GetInputByPort(address string, port int) (int, error) {
	command := make([]uint8, 4)

	//start with our instruction byte
	command[0] = 0x05

	//the SETUP bit - 0
	command[1] = 0x80

	//set to the 1 indexed port we care about
	command[2] = (uint8(port) + 1) | 0x80

	//we are controlling a single machine
	command[3] = 0x81
	for _, b := range command {
		log.Printf("%b", b)
	}

	resp, err := SendCommand(address, command)
	if err != nil {
		return 0, err
	}

	for _, b := range resp {
		log.Printf("%b", b)
	}

	output := resp[2] & 0x7F
	return int(output), nil
}
