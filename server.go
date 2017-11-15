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
	router := echo.New()
	router.Pre(middleware.RemoveTrailingSlash())
	router.Use(middleware.CORS())
	router.Use(CORS())

	// Use the `secure` routing group to require authentication
	secure := router.Group("", echo.WrapMiddleware(authmiddleware.Authenticate))

	// videoswitcher endpoints
	secure.GET("/:address/input/:input/:output", handlers.SwitchInput)
	secure.GET("/:address/front-lock/:bool", handlers.SetFrontLock)
	secure.GET("/:address/input/get/:port", handlers.GetInputByPort)

	// via endpoints
	secure.GET("/via/:address/reset", handlers.ResetVia)
	secure.GET("/via/:address/reboot", handlers.RebootVia)
	secure.GET("/via/:address/connected", handlers.GetViaConnectedStatus)

	port := ":8014"
	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	printHeader()
	router.StartServer(&server)
}

func CORS() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			return h(c)
		}
	}
}

func printHeader() {
	defer color.Unset()

	color.Set(color.FgHiYellow)
	fmt.Printf("\t\tKramer Microservice\n")

	// Videoswitcher
	fmt.Printf("Videoswitcher Endpoints:\n")

	color.Set(color.FgBlue)
	fmt.Printf("\t/:address/input/:input/:output\n")

	color.Set(color.FgHiCyan)
	fmt.Printf("\t\tChange the current input for a given output\n")

	color.Set(color.FgBlue)
	fmt.Printf("\t/:address/front-lock/:bool\n")

	color.Set(color.FgHiCyan)
	fmt.Printf("\t\tChange the front-button-lock status (true/false)\n")

	color.Set(color.FgBlue)
	fmt.Printf("\t/:address/input/get/:port\n")

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
