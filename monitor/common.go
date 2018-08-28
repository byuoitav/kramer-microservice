package monitor

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/byuoitav/common/structs"
	"github.com/fatih/color"
	"log"
	"net"
	"strings"
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

type Message struct {
	EventType string
	Action    string
	User      string `omitempty`
	State     string
}

// Old Send Command - Remove
/*
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
	if len(command.Command) > 0 {
		command.addAuth(false)
		command.writeCommand(conn)
	}

	return "", nil
}
*/

// Login to the VIA
func login(pconn *net.TCPConn) error {
	defer color.Unset()

	var cmd ViaCommand
	cmd.addAuth(true)
	cmd.Command = "Login"

	color.Set(color.FgBlue)
	log.Printf("Logging in...")

	b, err := xml.Marshal(cmd)
	if err != nil {
		return err
	}
	pconn.Write(b)

	color.Set(color.FgBlue)
	log.Printf("Login successful")

	return nil
}

func (c *ViaCommand) writeCommand(pconn *net.TCPConn) {
	defer color.Unset()

}

// Pass Authorization Creds
func (c *ViaCommand) addAuth(password bool) {
	c.Username = "su"
	if password {
		c.Password = "supass"
	}
}

// Get Connected to a VIA
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

// Read events and send them to console
func readPump(Pconn *net.TCPConn) {
	for {
		var m Message
		Buffer := make([]byte, 0, 2048)
		tmp := make([]byte, 256)

		r, err := Pconn.Read(tmp)
		if err != nil {
			err = errors.New(fmt.Sprintf("error reading from system: %s", err.Error()))
			log.Printf(err.Error())
		}
		//fmt.Printf("%s", r)
		Buffer = append(Buffer, tmp[:r]...)

		fmt.Println(string(Buffer))
		str := fmt.Sprintf("%s", Buffer)
		Out := strings.Split(str, "|")
		switch events := Out[0]; events {
		case "PList":
			m.EventType = Out[0]
			m.Action = Out[1]
			m.User = Out[2]
			m.State = Out[3]
		case "MediaStatus":
			m.EventType = Out[0]
			m.Action = Out[1]
			m.User = ""
			m.State = Out[2]
		case "DisplayStatus":
			m.EventType = Out[0]
			m.Action = Out[1]
			m.User = Out[2]
			m.State = Out[3]
		}
		/*
			if string.Contains(Out[0], "PList") {
				m.EventType = Out[0]
				m.Action = Out[1]
				m.User = Out[2]
				m.State = Out[3]
			}
			if string.Contains(Out0, "MediaStatus") {
				m.EventType = Out[0]
				m.Action = Out[1]
				m.User = ""
				m.State = Out[2]
			}
		*/
	}
}

func charInBuffer(toCheck byte, buffer []byte) bool {
	for _, b := range buffer {
		if toCheck == b {
			return true
		}
	}

	return false
}

// Build monitoring service for each VIA in a room
func StartMonitoring(device structs.Device) (Pconn *net.TCPConn) {
	//defer Pconn.Close()
	fmt.Printf("Building Connection and starting read buffer for %s\n", device.Address)
	addr := device.Address
	Pconn, err := PersistConnection(addr)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error reading response: %s\n", err.Error()))
	}

	go readPump(Pconn)
	//go writePump()

	return Pconn
}

// Build persistent connection with VIA
func PersistConnection(addr string) (*net.TCPConn, error) {
	defer color.Unset()
	color.Set(color.FgCyan)

	// get the connection
	log.Printf("Opening persistent telnet connection for reading events from %s", addr)
	Pconn, err := getConnection(addr)
	if err != nil {
		return nil, err
	}

	// login
	login(Pconn)

	return Pconn, nil
}
