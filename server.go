package main

import (
	"fmt"
	"net/http"

	"github.com/byuoitav/authmiddleware"
	"github.com/byuoitav/hateoas"
	"github.com/byuoitav/kramer-microservice/handlers"
	"github.com/fatih/color"
	"github.com/jessemillar/health"
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

	router.GET("/", echo.WrapHandler(http.HandlerFunc(hateoas.RootResponse)))
	router.GET("/health", echo.WrapHandler(http.HandlerFunc(health.Check)))

	//Functionality endpoints
	secure.GET("/:address/input/:input/:output", handlers.SwitchInput)
	secure.GET("/:address/front-lock/:bool", handlers.SetFrontLock)

	//Status endpoints
	//	secure.GET("/:address/input/map", handlers.GetCurrentInput)
	secure.GET("/:address/input/get/:port", handlers.GetInputByPort)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	color.Set(color.FgHiYellow)
	fmt.Printf("\t\tKramer Microservice\n")
	fmt.Printf("Endpoints:\n")
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
	color.Unset()

	router.StartServer(&server)
}
