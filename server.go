package main

import (
	//"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/byuoitav/authmiddleware"
	"github.com/byuoitav/common/db"
	"github.com/byuoitav/common/structs"
	"github.com/byuoitav/kramer-microservice/handlers"
	"github.com/byuoitav/kramer-microservice/handlers2000"
	"github.com/byuoitav/kramer-microservice/monitor"
	"github.com/byuoitav/kramer-microservice/videoswitcher"
	"github.com/fatih/color"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

/* global variable declaration */
// Changed: lowercase vars
var name string
var deviceList []structs.Device

func init() {
	name = os.Getenv("PI_HOSTNAME")
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
		log.Printf("%v\n", state)
		if (err != nil || state != "completed") && !(len(os.Getenv("DEV_ROUTER")) > 0 || len(os.Getenv("STOP_REPLICATION")) > 0) {
			log.Printf(color.RedString("Database replication in state %v. Retrying in 5 seconds.", state))
			time.Sleep(5 * time.Second)
			continue
		}
		log.Printf(color.GreenString("Database replication state: %v", state))

		devices, err := db.GetDB().GetDevicesByRoomAndRole(room, "EventRouter")
		if err != nil {
			log.Printf(color.RedString("Connecting to the Configuration DB failed, retrying in 5 seconds."))
			time.Sleep(5 * time.Second)
			continue
		}

		if len(devices) == 0 {
			//there's a chance that there ARE routers in the room, but the initial database replication is occuring.
			//we're good, keep going
			state, err := db.GetDB().GetStatus()
			if (err != nil || state != "completed") && !(len(os.Getenv("STOP_REPLICATION")) > 0) {
				log.Printf(color.RedString("Database replication in state %v. Retrying in 5 seconds.", state))
				time.Sleep(5 * time.Second)
				continue
			}
		}
		break
		log.Printf(color.BlueString("Connection to the Configuration DB established."))
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
	router := echo.New()
	router.Pre(middleware.RemoveTrailingSlash())
	router.Use(middleware.CORS())

	// Use the `secure` routing group to require authentication
	secure := router.Group("", echo.WrapMiddleware(authmiddleware.Authenticate))

	//start the VIA monitoring connection if the Controller is CP1
	if strings.Contains(name, "-CP1") {
		for _, device := range deviceList {
			go monitor.StartMonitoring(device)
		}
	}

	// videoswitcher endpoints
	secure.GET("/:address/welcome/:bool/input/:input/:output", handlers.SwitchInput)
	secure.GET("/:address/welcome/:bool/front-lock/:bool2", handlers.SetFrontLock)
	secure.GET("/:address/welcome/:bool/input/get/:port", handlers.GetInputByPort)

	secure.GET("/2000/:address/input/:input/:output", handlers2000.SwitchInput)
	secure.GET("/2000/:address/input/get/:port", handlers2000.GetInputByPort)

	// via functionality endpoints
	secure.GET("/via/:address/reset", handlers.ResetVia)
	secure.GET("/via/:address/reboot", handlers.RebootVia)

	// Set the volume
	secure.GET("/via/:address/volume/set/:volvalue", handlers.SetViaVolume)

	// via informational endpoints
	secure.GET("/via/:address/connected", handlers.GetViaConnectedStatus)
	secure.GET("/via/:address/volume/level", handlers.GetViaVolume)

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
}
