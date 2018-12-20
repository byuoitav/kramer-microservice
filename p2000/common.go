package p2000

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/byuoitav/common/log"

	"github.com/fatih/color"
)

func SendCommand(address string, command []uint8) ([]uint8, error) {
	log.L.Infof("sending %X", command)

	conn, err := getConnection(address)
	if err != nil {
		log.L.Infof(color.HiRedString(err.Error()))
		return []uint8{}, err
	}

	defer conn.Close()

	//we just write the bytes
	conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
	num, err := conn.Write(command)
	if err != nil {
		log.L.Infof(color.HiRedString(err.Error()))
		return []uint8{}, err
	}
	if num != 4 {
		msg := fmt.Sprintf("There were an invalid number of bytes written, should have been 4, received %v", num)
		log.L.Infof(color.HiRedString(msg))
		return []uint8{}, errors.New(msg)
	}

	//get four bytes
	resp := make([]uint8, 4)

	log.L.Infof("Written. Reading.")

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	num, err = conn.Read(resp)
	if err != nil {
		log.L.Infof(color.HiRedString(err.Error()))
		return []uint8{}, err
	}
	if num != 4 {
		msg := fmt.Sprintf("There were an invalid number of bytes read, should have been 4, received %v", num)
		log.L.Infof(color.HiRedString(msg))
		return []uint8{}, errors.New(msg)
	}
	log.L.Infof("Recieved %+X", resp)

	return resp, nil
}

func getConnection(address string) (*net.TCPConn, error) {
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
