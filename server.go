package main

import (
	//"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/byuoitav/common"
	"github.com/byuoitav/common/db"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/structs"
	"github.com/byuoitav/common/v2/auth"
	"github.com/byuoitav/kramer-microservice/handlers"
	"github.com/byuoitav/kramer-microservice/handlers2000"
	"github.com/byuoitav/kramer-microservice/monitor"
	"github.com/byuoitav/kramer-microservice/videoswitcher"
	"github.com/fatih/color"
)

/* global variable declaration */
// Changed: lowercase vars
var name string
var deviceList []structs.Device

func init() {

	if len(os.Getenv("ROOM_SYSTEM")) == 0 {
		log.L.Debugf("System is not tied to a specific room. Will not start via monitoring")
		return
	}

	name = os.Getenv("SYSTEM_ID")
	var err error
	fmt.Printf("Gathering information for %s from database\n", name)

	s := strings.Split(name, "-")
	sa := s[0:2]
	room := strings.Join(sa, "-")
	//fmt.Printf("Waiting for database entry for %s\n", name)
	fmt.Printf("Waiting for database . . . .\n")
	for {
		// Pull room information from db
		state, err := db.GetDB().GetStatus()
		log.L.Debugf("%v\n", state)
		//+deploy not_requried
		if (err != nil || state != "completed") && !(len(os.Getenv("DEV_ROUTER")) > 0 || len(os.Getenv("STOP_REPLICATION")) > 0) {
			log.L.Debugf(color.RedString("Database replication in state %v. Retrying in 5 seconds.", state))
			time.Sleep(5 * time.Second)
			continue
		}
		log.L.Debugf(color.GreenString("Database replication state: %v", state))

		devices, err := db.GetDB().GetDevicesByRoomAndRole(room, "EventRouter")
		if err != nil {
			log.L.Debugf(color.RedString("Connecting to the Configuration DB failed, retrying in 5 seconds."))
			time.Sleep(5 * time.Second)
			continue
		}

		if len(devices) == 0 {
			//there's a chance that there ARE routers in the room, but the initial database replication is occuring.
			//we're good, keep going
			state, err := db.GetDB().GetStatus()
			if (err != nil || state != "completed") && !(len(os.Getenv("STOP_REPLICATION")) > 0) {
				log.L.Debugf(color.RedString("Database replication in state %v. Retrying in 5 seconds.", state))
				time.Sleep(5 * time.Second)
				continue
			}
		}
		log.L.Debugf(color.BlueString("Connection to the Configuration DB established."))
		break
	}
	deviceList, err = db.GetDB().GetDevicesByRoomAndType(room, "via-connect-pro")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func main() {

	//start the router
	go videoswitcher.StartRouter()

	port := ":8014"
	router := common.NewRouter()

	write := router.Group("", auth.AuthorizeRequest("write-state", "room", auth.LookupResourceFromAddress))
	read := router.Group("", auth.AuthorizeRequest("read-state", "room", auth.LookupResourceFromAddress))

	//start the VIA monitoring connection if the Controller is CP1
	if strings.Contains(name, "-CP1") && len(os.Getenv("ROOM_SYSTEM")) > 0 {
		for _, device := range deviceList {
			go monitor.StartMonitoring(device)
		}
	}

	// videoswitcher endpoints
	read.GET("/:address/welcome/:bool/input/:input/:output", handlers.SwitchInput)
	read.GET("/:address/welcome/:bool/front-lock/:bool2", handlers.SetFrontLock)
	read.GET("/:address/welcome/:bool/input/get/:port", handlers.GetInputByPort)

	write.GET("/2000/:address/input/:input/:output", handlers2000.SwitchInput)
	read.GET("/2000/:address/input/get/:port", handlers2000.GetInputByPort)

	// via functionality endpoints
	write.GET("/via/:address/reset", handlers.ResetVia)
	write.GET("/via/:address/reboot", handlers.RebootVia)

	// Set the volume
	write.GET("/via/:address/volume/set/:volvalue", handlers.SetViaVolume)

	// via informational endpoints
	read.GET("/via/:address/connected", handlers.GetViaConnectedStatus)
	read.GET("/via/:address/volume/level", handlers.GetViaVolume)
	read.GET("/via/:address/hardware", handlers.GetViaHardwareInfo)
	read.GET("/via/:address/users/status", handlers.GetStatusOfUsers)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	printHeader()
	router.StartServer(&server)
}

func printHeader() {
	defer color.Unset()

	color.Set(color.FgHiYellow)
	fmt.Printf("\t\tKramer Microservice\n")

	// Videoswitcher
	fmt.Printf("Videoswitcher Endpoints:\n")

	color.Set(color.FgBlue)
	fmt.Printf("\t/:address/welcome/:bool/input/:input/:output\n")

	color.Set(color.FgHiCyan)
	fmt.Printf("\t\tChange the current input for a given output\n")

	color.Set(color.FgBlue)
	fmt.Printf("\t/:address/welcome/:bool/front-lock/:bool2\n")

	color.Set(color.FgHiCyan)
	fmt.Printf("\t\tChange the front-button-lock status (true/false)\n")

	color.Set(color.FgBlue)
	fmt.Printf("\t/:address/welcome/:bool/input/get/:port\n")

	color.Set(color.FgHiCyan)
	fmt.Printf("\t\tGet the current input for a given output port\n")

	// VIA
	color.Set(color.FgHiYellow)
	fmt.Printf("VIA Endpoints:\n")

	color.Set(color.FgBlue)
	fmt.Printf("\t/via/:address/reboot\n")

	color.Set(color.FgHiCyan)
	fmt.Printf("\t\tReboot a VIA\n")

	color.Set(color.FgBlue)
	fmt.Printf("\t/via/:address/reset\n")

	color.Set(color.FgHiCyan)
	fmt.Printf("\t\tReset a VIA\n")

	color.Set(color.FgBlue)
	fmt.Printf("\t/via/:address/connected\n")

	color.Set(color.FgHiCyan)
	fmt.Printf("\t\tGet connected status of a Via device\n")

	color.Set(color.FgBlue)
	fmt.Printf("\t/via/:address/volume/set/:volvalue\n")

	color.Set(color.FgHiCyan)
	fmt.Printf("\t\tSet volume level on a VIA device\n")

	color.Set(color.FgBlue)
	fmt.Printf("\t/via/:address/volume/level\n")

	color.Set(color.FgHiCyan)
	fmt.Printf("\t\tGet volume level on a VIA device\n")

	color.Set(color.FgBlue)
	fmt.Printf("\t/via/:address/hardware\n")

	color.Set(color.FgHiCyan)
	fmt.Printf("\t\tGet the hardware information of a VIA device\n")

	color.Set(color.FgBlue)
	fmt.Printf("\t/via/:address/users/status\n")

	color.Set(color.FgHiCyan)
	fmt.Printf("\t\tGet the status of users logged into a VIA device\n")
}
