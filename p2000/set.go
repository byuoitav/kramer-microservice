package p2000

//returns the input
func SetOutput(address string, input, output int) (int, error) {
	command := make([]uint8, 4)

	command[0] = 0x01
	command[1] = 0x80 | (uint8(input) + 1)
	command[2] = 0x80 | (uint8(output) + 1)
	command[3] = 0x81

	resp, err := SendCommand(address, command)
	if err != nil {
		return 0, err
	}

	toReturn := resp[1] & 0x7F
	return int(toReturn) - 1, nil

}
