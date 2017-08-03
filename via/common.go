package via

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/fatih/color"
)

type ViaCommand struct {
	XMLName  xml.Name `xml:"P"`
	Username string   `xml:"UN"`
	Password string   `xml:"Pwd"`
	Command  string   `xml:"Cmd"`
	Param1   string   `xml:"P1,omitempty"`
	Param2   string   `xml:"P2,omitempty"`
	Param3   string   `xml:"P3,omitempty"`
	Param4   string   `xml:"P4,omitempty"`
	Param5   string   `xml:"P5,omitempty"`
	Param6   string   `xml:"P6,omitempty"`
	Param7   string   `xml:"P7,omitempty"`
	Param8   string   `xml:"P8,omitempty"`
	Param9   string   `xml:"P9,omitempty"`
	Param10  string   `xml:"P10,omitempty"`
}

func SendCommand(command ViaCommand, addr string) (string, error) {
	defer color.Unset()
	color.Set(color.FgCyan)

	// get the connection
	log.Printf("Opening telnet connection with %s", addr)
	conn, err := getConnection(addr)
	if err != nil {
		return "", err
	}

	// login
	login(conn)

	// write command
	command.addAuth(false)
	command.writeCommand(conn)

	// get response
	resp, err := readUntil('\n', conn, 5)
	if len(string(resp)) > 0 {
		color.Set(color.FgBlue)
		log.Printf("Response from device: %s", resp)
	}

	return string(resp), nil
}

func login(conn *net.TCPConn) error {
	defer color.Unset()

	var cmd ViaCommand
	cmd.addAuth(true)
	cmd.Command = "Login"

	color.Set(color.FgBlue)
	log.Printf("Logging in...")

	err := cmd.writeCommand(conn)
	if err != nil {
		return err
	}

	color.Set(color.FgBlue)
	log.Printf("Login successful")

	return nil
}

func (c *ViaCommand) writeCommand(conn *net.TCPConn) error {
	defer color.Unset()

	// read welcome message
	_, err := readUntil('\n', conn, 3)
	if err != nil {
		return err
	}

	b, err := xml.Marshal(c)
	if err != nil {
		return err
	}

	color.Set(color.FgMagenta)
	if len(c.Password) == 0 {
		log.Printf("Sending command: %s", b)
		color.Set(color.FgMagenta)
	}

	conn.Write(b)
	return nil
}

func (c *ViaCommand) addAuth(password bool) {
	c.Username = "su"
	if password {
		c.Password = "supass"
	}
}

func getConnection(address string) (*net.TCPConn, error) {
	radder, err := net.ResolveTCPAddr("tcp", address+":9982")
	if err != nil {
		err = errors.New(fmt.Sprintf("error resolving address : %s", err.Error()))
		log.Printf(err.Error())
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, radder)
	if err != nil {
		err = errors.New(fmt.Sprintf("error dialing address : %s", err.Error()))
		log.Printf(err.Error())
		return nil, err
	}

	return conn, nil
}

func readUntil(delimeter byte, conn *net.TCPConn, timeoutInSeconds int) ([]byte, error) {
	conn.SetReadDeadline(time.Now().Add(time.Duration(int64(timeoutInSeconds)) * time.Second))

	buffer := make([]byte, 128)
	message := []byte{}

	for !charInBuffer(delimeter, buffer) {
		_, err := conn.Read(buffer)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error reading response: %s", err.Error()))
			color.Set(color.FgRed)
			log.Printf("%s", err.Error())
			color.Unset()
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
