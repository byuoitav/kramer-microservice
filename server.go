package main

import (
	"fmt"
	"net/http"

	"github.com/byuoitav/authmiddleware"
	"github.com/byuoitav/kramer-microservice/handlers"
	"github.com/fatih/color"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	port := ":8014"
	router := echo.New()
	router.Pre(middleware.RemoveTrailingSlash())
	router.Use(middleware.CORS())

	// Use the `secure` routing group to require authentication
	secure := router.Group("", echo.WrapMiddleware(authmiddleware.Authenticate))

	// videoswitcher endpoints
	secure.GET("/:address/welcome/:bool/input/:input/:output", handlers.SwitchInput)
	secure.GET("/:address/welcome/:bool/front-lock/:bool2", handlers.SetFrontLock)
	secure.GET("/:address/welcome/:bool/input/get/:port", handlers.GetInputByPort)

	// via endpoints
	secure.GET("/via/:address/reset", handlers.ResetVia)
	secure.GET("/via/:address/reboot", handlers.RebootVia)
	secure.GET("/via/:address/connected", handlers.GetViaConnectedStatus)

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
	fmt.Printf("\t\tGet connected status of a via\n")
}
