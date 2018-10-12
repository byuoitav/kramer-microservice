package monitor

import (
	"encoding/xml"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/byuoitav/common/events"
	"github.com/byuoitav/common/structs"
	"github.com/byuoitav/kramer-microservice/via"
	"github.com/fatih/color"
)

const (
	// Intervals to wait between retry attempts
	reconnInterval = 10 * time.Second

	// Ping Internal (in milliseconds, because it cares)
	pingInterval = 60000
)

var (
	pihost     string
	hostname   string
	buildingID string
	room       string
)

func init() {
	var err error

	pihost = os.Getenv("PI_HOSTNAME")
	if len(pihost) == 0 {
		log.Fatalf("PI_HOSTNAME not set.")
	}

	hostname, err = os.Hostname()
	if err != nil {
		hostname = pihost
	}

	split := strings.Split(pihost, "-")
	buildingID = split[0]
	room = split[1]
}

type message struct {
	EventType string
	Action    string
	User      string
	State     string
}

// Ping over connection to keep alive.
func pingTest(pconn *net.TCPConn) error {
	defer color.Unset()
	color.Set(color.FgCyan)
	var c via.Command
	c.Username = "su"
	c.Command = "IpInfo"
	log.Printf("Oh ho, Pongo, you old rascal!")
	b, err := xml.Marshal(c)
	if err != nil {
		return err
	}
	_, err = pconn.Write(b)
	if err != nil {
		return err
	}
	return err
}

// Retry connection if connection has failed
func retryViaConnection(device structs.Device, pconn *net.TCPConn, event events.Event) {
	log.Printf(color.HiMagentaString("[retry] Retrying Connection to VIA"))
	addr := device.Address
	pconn, err := via.PersistConnection(addr)
	for err != nil {
		log.Printf(color.RedString("Retry Failed, Trying again in 10 seconds"))
		time.Sleep(reconnInterval)
		pconn, err = via.PersistConnection(addr)
	}

	go readPump(device, pconn, event)
	go writePump(device, pconn)
}

// Read events and send them to console
func readPump(device structs.Device, pconn *net.TCPConn, event events.Event) {
	// defer closing connection
	defer func(device structs.Device) {
		pconn.Close()
		log.Printf(color.HiRedString("Connection to VIA %v is dying.", device.Address))
		log.Printf(color.HiRedString("Trying to reconnect........"))
		//retry connection to VIA device
		retryViaConnection(device, pconn, event)
	}(device)
	timeoutDuration := 300 * time.Second

	for {
		var m message
		//Set buffers for accepting data
		Buffer := make([]byte, 0, 2048)
		tmp := make([]byte, 256)

		//set deadline for reads - keep the connection alive during that time
		pconn.SetReadDeadline(time.Now().Add(timeoutDuration))
		//start reader to read into buffer
		r, err := pconn.Read(tmp)
		if err != nil {
			err = fmt.Errorf("error reading from system: %s", err.Error())
			log.Printf(err.Error())
			return
		}
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
		event.Timestamp = time.Now().Format(time.RFC3339)
		event.Event.EventInfoKey = m.EventType
		event.Event.EventInfoValue = m.State
		event.Event.User = m.User

		// changed: add event stuff
		eventNode().PublishEvent(events.Metrics, event)
	}
}

func writePump(device structs.Device, pconn *net.TCPConn) {
	// defer closing connection
	defer func(device structs.Device) {
		pconn.Close()
		log.Printf(color.HiRedString("Error on write pump for %v. Write pump closing.", device.Address))
	}(device)
	ticker := time.NewTicker(pingInterval * time.Millisecond)
	// Once the pingInterval is reached, execute the ping -
	// On Error, return and execute deferred to close the connection
	for range ticker.C {
		err := pingTest(pconn)
		if err != nil {
			log.Printf(color.HiRedString("Ping Failed Error: %v", err))
			return
		}
	}
}

// StartMonitoring service for each VIA in a room
func StartMonitoring(device structs.Device) *net.TCPConn {
	fmt.Printf("Building Connection and starting read buffer for %s\n", device.Address)
	addr := device.Address
	pconn, err := via.PersistConnection(addr)
	if err != nil {
		err = fmt.Errorf("error reading response: %s", err.Error())
		return nil
	}

	// start event node
	_ = eventNode()

	// build base event to send along with each event
	event := events.Event{
		Hostname:         pihost,
		LocalEnvironment: true,
		Building:         buildingID,
		Room:             room,
		Event: events.EventInfo{
			Type:       events.DETAILSTATE,
			Requestor:  hostname,
			EventCause: events.AUTOGENERATED,
			Device:     device.Name,
			DeviceID:   device.ID,
		},
	}

	go readPump(device, pconn, event)
	go writePump(device, pconn)
	return pconn
}

var once sync.Once
var node *events.EventNode

func eventNode() *events.EventNode {
	once.Do(func() {
		router := os.Getenv("EVENT_ROUTER_ADDRESS")
		if len(router) == 0 {
			log.Fatalf("EVENT_ROUTER_ADDRESS is not set.")
		}
		node = events.NewEventNode("Kramer Microservice", router, []string{})
	})

	return node
}
