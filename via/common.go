package via

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net"
)

type ViaCommand struct {
	XMLName  xml.Name `xml:"P"`
	Username string   `xml:"UN"`
	Password string   `xml:"Pwd"`
	Command  string   `xml:"Cmd"`
	Param1   string   `xml:"P1"`
	Param2   string   `xml:"P2"`
	Param3   string   `xml:"P3"`
	Param4   string   `xml:"P4"`
	Param5   string   `xml:"P5"`
	Param6   string   `xml:"P6"`
	Param7   string   `xml:"P7"`
	Param8   string   `xml:"P8"`
	Param9   string   `xml:"P9"`
	Param10  string   `xml:"P10"`
}

func sendCommand(command string, addr string, username string, password string) error {
	log.Printf("Sending command %v to %v", command, addr)

	//get the connection
	_, err := getConnection(addr)
	if err != nil {
		log.Printf("There was a problem getting the connection")
		return err
	}
	return nil
}

func login(conn *net.TCPConn, username string, password string) error {
	//successMessage := "Login Successful."

	log.Printf("Logging in")
	return nil

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
