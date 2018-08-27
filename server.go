package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

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
var Name string
var DeviceList []structs.Device

func init() {
	Name = os.Getenv("PI_HOSTNAME")
	var err error
	fmt.Printf("Gathering information for %s from database\n", Name)

	s := strings.Split(Name, "-")
	sa := s[0:2]
	room := strings.Join(sa, "-")
	fmt.Printf("Waiting for database entry for %s\n", Name)
	// Pull room information from db
	DeviceList, err = db.GetDB().GetDevicesByRoomAndType(room, "via-connect-pro")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	fmt.Printf("Device List: \n")
	r, _ := json.Marshal(DeviceList)
	fmt.Printf(string(r))
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
	//hostname := os.Getenv("PI_HOSTNAME")
	if strings.Contains(Name, "-CP1") {
		for _, Device := range DeviceList {
			go monitor.StartMonitoring(Device)
			resp, _ := json.Marshal(Device)
			fmt.Printf(string(resp))
			fmt.Printf("\n")
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
