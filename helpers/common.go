package helpers

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/fatih/color"
)

const CARRIAGE_RETURN = 0x0D
const LINE_FEED = 0x0A
const SPACE = 0x20

// takes a command and sends it to the device
// should it also return what the response of the query was?
func SendCommand(address string, command string) (string, error) {
	defer color.Unset()
	color.Set(color.FgCyan)

	// open telnet connection with address
	log.Printf("Opening telnet connection with %s", address)
	conn, err := GetConnection(address)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// read the welcom message
	_, err = readUntil(CARRIAGE_RETURN, conn, 3)

	// write command
	log.Printf("Sending command %s", command)
	command += string(CARRIAGE_RETURN) + string(LINE_FEED)
	conn.Write([]byte(command))

	// get response
	resp, err := readUntil(LINE_FEED, conn, 5)
	color.Set(color.FgBlue)
	log.Printf("Response from device: %s", resp)

	return string(resp), nil
}

func GetConnection(address string) (*net.TCPConn, error) {
	addr, err := net.ResolveTCPAddr("tcp", address+":5000")
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}

	return conn, err
}

func readUntil(delimeter byte, conn *net.TCPConn, timeoutInSeconds int) ([]byte, error) {

	conn.SetReadDeadline(time.Now().Add(time.Duration(int64(timeoutInSeconds)) * time.Second))

	buffer := make([]byte, 128)
	message := []byte{}

	for !charInBuffer(delimeter, buffer) {
		_, err := conn.Read(buffer)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error reading response: %s", err.Error()))
			log.Printf("%s", err.Error())
			return message, err
		}

		message = append(message, buffer...)
	}
	return bytes.Trim(message, "\x00"), nil
}

func charInBuffer(toCheck byte, buffer []byte) bool {
	for _, b := range buffer {
		if toCheck == b {
			return true
		}
	}

	return false
}
